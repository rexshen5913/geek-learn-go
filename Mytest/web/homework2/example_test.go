package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWebFrameworkBasicUsage 展示 Web 框架的基本使用方法
func TestWebFrameworkBasicUsage(t *testing.T) {
	// 創建 HTTP 服務器
	server := NewHttpServer()

	// 註冊基本路由
	server.Get("/hello", func(ctx *Context) {
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("Hello, World!")
	})

	// 註冊帶參數的路由
	server.Get("/users/:id", func(ctx *Context) {
		id, err := ctx.PathValue("id")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("Invalid user ID")
			return
		}
		ctx.RespJSON(http.StatusOK, map[string]string{
			"message": "User found",
			"id":      id,
		})
	})

	// 註冊 POST 路由處理 JSON
	server.Post("/users", func(ctx *Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := ctx.BindJSON(&user); err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("Invalid JSON")
			return
		}

		// 模擬創建用戶
		response := map[string]interface{}{
			"id":      123,
			"name":    user.Name,
			"email":   user.Email,
			"message": "User created successfully",
		}

		ctx.RespJSON(http.StatusCreated, response)
	})

	// 測試基本路由
	t.Run("Basic Hello Route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if w.Body.String() != "Hello, World!" {
			t.Errorf("expected 'Hello, World!', got '%s'", w.Body.String())
		}
	})

	// 測試路徑參數
	t.Run("Path Parameter Route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/456", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to parse JSON response: %v", err)
		}

		if response["id"] != "456" {
			t.Errorf("expected id '456', got '%s'", response["id"])
		}

		if response["message"] != "User found" {
			t.Errorf("expected message 'User found', got '%s'", response["message"])
		}
	})

	// 測試 JSON 請求處理
	t.Run("JSON Request Handling", func(t *testing.T) {
		userData := map[string]string{
			"name":  "張三",
			"email": "zhangsan@example.com",
		}
		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to parse JSON response: %v", err)
		}

		if response["name"] != "張三" {
			t.Errorf("expected name '張三', got '%v'", response["name"])
		}

		if response["email"] != "zhangsan@example.com" {
			t.Errorf("expected email 'zhangsan@example.com', got '%v'", response["email"])
		}

		if response["message"] != "User created successfully" {
			t.Errorf("expected message 'User created successfully', got '%v'", response["message"])
		}
	})
}

// TestMiddlewareUsage 展示中間件的使用方法
func TestMiddlewareUsage(t *testing.T) {
	var logs []string

	// 日誌中間件
	loggingMiddleware := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			logs = append(logs, "Request: "+ctx.Req.Method+" "+ctx.Req.URL.Path)
			next(ctx)
			logs = append(logs, "Response sent")
		}
	}

	// 創建帶中間件的服務器
	server := NewHttpServer(ServerWithMiddleware(loggingMiddleware))

	server.Get("/api/test", func(ctx *Context) {
		logs = append(logs, "Handler executed")
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("API response")
	})

	logs = []string{} // 重置日誌

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	// 驗證中間件執行順序
	expectedLogs := []string{
		"Request: GET /api/test",
		"Handler executed",
		"Response sent",
	}

	if len(logs) != len(expectedLogs) {
		t.Fatalf("expected %d logs, got %d: %v", len(expectedLogs), len(logs), logs)
	}

	for i, expected := range expectedLogs {
		if logs[i] != expected {
			t.Errorf("log[%d]: expected '%s', got '%s'", i, expected, logs[i])
		}
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "API response" {
		t.Errorf("expected 'API response', got '%s'", w.Body.String())
	}
}

// TestQueryParameterHandling 展示查詢參數處理
func TestQueryParameterHandling(t *testing.T) {
	server := NewHttpServer()

	server.Get("/search", func(ctx *Context) {
		// 獲取必需的查詢參數
		query, err := ctx.QueryValue("q")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("Missing search query")
			return
		}

		// 獲取可選的查詢參數，帶默認值
		limit := "10"
		if l, err := ctx.QueryValue("limit"); err == nil {
			limit = l
		}

		category := "all"
		if c, err := ctx.QueryValue("category"); err == nil {
			category = c
		}

		response := map[string]string{
			"query":    query,
			"limit":    limit,
			"category": category,
			"status":   "success",
		}

		ctx.RespJSON(http.StatusOK, response)
	})

	// 測試完整的查詢參數
	t.Run("Full Query Parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/search?q=golang&limit=20&category=programming", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["query"] != "golang" {
			t.Errorf("expected query 'golang', got '%s'", response["query"])
		}

		if response["limit"] != "20" {
			t.Errorf("expected limit '20', got '%s'", response["limit"])
		}

		if response["category"] != "programming" {
			t.Errorf("expected category 'programming', got '%s'", response["category"])
		}
	})

	// 測試部分查詢參數（使用默認值）
	t.Run("Partial Query Parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/search?q=python", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["query"] != "python" {
			t.Errorf("expected query 'python', got '%s'", response["query"])
		}

		if response["limit"] != "10" {
			t.Errorf("expected default limit '10', got '%s'", response["limit"])
		}

		if response["category"] != "all" {
			t.Errorf("expected default category 'all', got '%s'", response["category"])
		}
	})

	// 測試缺少必需參數
	t.Run("Missing Required Parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/search", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}

		if w.Body.String() != "Missing search query" {
			t.Errorf("expected error message, got '%s'", w.Body.String())
		}
	})
}
