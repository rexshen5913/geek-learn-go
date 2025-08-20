package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHttpServer_RouteRegistration 測試路由註冊功能
func TestHttpServer_RouteRegistration(t *testing.T) {
	server := NewHttpServer()

	// 測試基本路由註冊
	server.Get("/users", func(ctx *Context) {
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("users list")
	})

	server.Post("/users", func(ctx *Context) {
		ctx.RespStatusCode = http.StatusCreated
		ctx.RespData = []byte("user created")
	})

	// 測試路徑參數路由
	server.Get("/users/:id", func(ctx *Context) {
		id, err := ctx.PathValue("id")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("missing id")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte(fmt.Sprintf("user %s", id))
	})

	// 測試通配符路由
	server.Get("/static/*", func(ctx *Context) {
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("static file")
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /users",
			method:         http.MethodGet,
			path:           "/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "users list",
		},
		{
			name:           "POST /users",
			method:         http.MethodPost,
			path:           "/users",
			expectedStatus: http.StatusCreated,
			expectedBody:   "user created",
		},
		{
			name:           "GET /users/123",
			method:         http.MethodGet,
			path:           "/users/123",
			expectedStatus: http.StatusOK,
			expectedBody:   "user 123",
		},
		{
			name:           "GET /static/css/main.css",
			method:         http.MethodGet,
			path:           "/static/css/main.css",
			expectedStatus: http.StatusOK,
			expectedBody:   "static file",
		},
		{
			name:           "GET /notfound",
			method:         http.MethodGet,
			path:           "/notfound",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tt.expectedBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

// TestHttpServer_Middleware 測試中間件功能
func TestHttpServer_Middleware(t *testing.T) {
	var executed []string

	middleware1 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executed = append(executed, "middleware1-before")
			next(ctx)
			executed = append(executed, "middleware1-after")
		}
	}

	middleware2 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executed = append(executed, "middleware2-before")
			next(ctx)
			executed = append(executed, "middleware2-after")
		}
	}

	server := NewHttpServer(ServerWithMiddleware(middleware1, middleware2))

	server.Get("/test", func(ctx *Context) {
		executed = append(executed, "handler")
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("test response")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executed) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executed))
	}

	for i, expected := range expectedOrder {
		if executed[i] != expected {
			t.Errorf("execution order[%d]: expected '%s', got '%s'", i, expected, executed[i])
		}
	}
}

// TestHttpServer_JSON 測試 JSON 處理功能
func TestHttpServer_JSON(t *testing.T) {
	server := NewHttpServer()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// 測試 JSON 響應
	server.Get("/users/json", func(ctx *Context) {
		user := User{ID: 1, Name: "張三"}
		ctx.RespJSON(http.StatusOK, user)
	})

	// 測試 JSON 請求解析
	server.Post("/users/json", func(ctx *Context) {
		var user User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("invalid json")
			return
		}

		// 返回接收到的用戶信息
		ctx.RespJSON(http.StatusCreated, user)
	})

	t.Run("GET JSON response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/json", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
		}

		var user User
		if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if user.ID != 1 || user.Name != "張三" {
			t.Errorf("unexpected user data: %+v", user)
		}
	})

	t.Run("POST JSON request", func(t *testing.T) {
		user := User{ID: 2, Name: "李四"}
		jsonData, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/users/json", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var responseUser User
		if err := json.Unmarshal(w.Body.Bytes(), &responseUser); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if responseUser.ID != user.ID || responseUser.Name != user.Name {
			t.Errorf("expected user %+v, got %+v", user, responseUser)
		}
	})
}

// TestHttpServer_QueryParams 測試查詢參數功能
func TestHttpServer_QueryParams(t *testing.T) {
	server := NewHttpServer()

	server.Get("/search", func(ctx *Context) {
		query, err := ctx.QueryValue("q")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("missing query parameter")
			return
		}

		limit := "10" // 默認值
		if l, err := ctx.QueryValue("limit"); err == nil {
			limit = l
		}

		response := fmt.Sprintf("searching for '%s' with limit %s", query, limit)
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte(response)
	})

	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "with query and limit",
			url:          "/search?q=golang&limit=20",
			expectedCode: http.StatusOK,
			expectedBody: "searching for 'golang' with limit 20",
		},
		{
			name:         "with query only",
			url:          "/search?q=python",
			expectedCode: http.StatusOK,
			expectedBody: "searching for 'python' with limit 10",
		},
		{
			name:         "missing query parameter",
			url:          "/search",
			expectedCode: http.StatusBadRequest,
			expectedBody: "missing query parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tt.expectedBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

// TestHttpServer_RouteMiddleware 測試路由級別的中間件
func TestHttpServer_RouteMiddleware(t *testing.T) {
	server := NewHttpServer()

	var executed []string

	// 路由級別的中間件
	authMiddleware := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executed = append(executed, "auth-check")
			// 模擬驗證邏輯
			if ctx.Req.Header.Get("Authorization") == "" {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("unauthorized")
				return
			}
			next(ctx)
		}
	}

	logMiddleware := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executed = append(executed, "log-request")
			next(ctx)
			executed = append(executed, "log-response")
		}
	}

	// 註冊路由中間件
	server.Use(http.MethodGet, "/protected", authMiddleware, logMiddleware)

	// 註冊處理函數
	server.Get("/protected", func(ctx *Context) {
		executed = append(executed, "protected-handler")
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("protected resource")
	})

	t.Run("unauthorized access", func(t *testing.T) {
		executed = []string{} // 重置

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		// 根據中間件的註冊順序，authMiddleware 是第一個，logMiddleware 是第二個
		// 但執行時是從後往前組裝鏈條，所以實際執行順序是後註冊的先執行外層
		expectedExecution := []string{"auth-check"}
		if len(executed) != len(expectedExecution) {
			t.Fatalf("expected %d executions, got %d: %v", len(expectedExecution), len(executed), executed)
		}

		for i, expected := range expectedExecution {
			if executed[i] != expected {
				t.Errorf("execution[%d]: expected '%s', got '%s'", i, expected, executed[i])
			}
		}
	})

	t.Run("authorized access", func(t *testing.T) {
		executed = []string{} // 重置

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// 根據中間件執行順序：logMiddleware 在外層，authMiddleware 在內層
		expectedExecution := []string{
			"auth-check",
			"log-request",
			"protected-handler",
			"log-response",
		}

		if len(executed) != len(expectedExecution) {
			t.Fatalf("expected %d executions, got %d: %v", len(expectedExecution), len(executed), executed)
		}

		for i, expected := range expectedExecution {
			if executed[i] != expected {
				t.Errorf("execution[%d]: expected '%s', got '%s'", i, expected, executed[i])
			}
		}
	})
	server.Start(":8080")
}

// TestHttpServer_ServerInterface 測試 Server 接口實現
func TestHttpServer_ServerInterface(t *testing.T) {
	// 測試 HttpServer 是否正確實現了 Server 接口
	var _ Server = &HttpServer{}

	server := NewHttpServer()

	// 測試接口方法是否可用
	server.addRoute(http.MethodGet, "/test", func(ctx *Context) {
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("interface test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if strings.TrimSpace(w.Body.String()) != "interface test" {
		t.Errorf("expected body 'interface test', got '%s'", strings.TrimSpace(w.Body.String()))
	}
}
