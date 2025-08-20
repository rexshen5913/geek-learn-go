package web

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/likexian/gokit/assert"
)

func TestRouter_AddRoute(t *testing.T) {
	// 第一個步驟是構造路由樹
	// 第二個步驟是驗證路由樹
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		// {
		// 	method: http.MethodGet,
		// 	path:   "/*",
		// },
		// {
		// 	method: http.MethodGet,
		// 	path:   "/*/*",
		// },
		// {
		// 	method: http.MethodGet,
		// 	path:   "/*/abc",
		// },
		// {
		// 	method: http.MethodGet,
		// 	path:   "/*/abc/*",
		// },
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 這裡斷言路由樹與預期得一模一樣
	// /user/home 會是以 / 為根節點，user 為子節點，home 為 user 的子節點
	// 在 home 那邊會有 handler 的註冊
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
						starChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := r.equal(wantRouter)
	assert.True(t, ok, msg)

	// 斷言兩者相等
	// 這裡無法使用，是因為有 handleFunc 的關係
	// 因爲那是無法比較的，func 只能比較地址
	// 所以這裡需要自己實現一個 Equal 函數
	// assert.Equal(t, wantRouter, r)

	r = newRouter()
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "a/b/c", mockHandler)
	})
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	})
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "//a/b/c", mockHandler)
	})
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/a//b/c", mockHandler)
	})
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/a/b//c", mockHandler)
	})
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/a/b///c", mockHandler)
	})

	// 測試重複路徑
	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})
	r.addRoute(http.MethodGet, "/user", mockHandler)
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/user", mockHandler)
	})

	// 測試 * 和 : 路由不能同時註冊
	r = newRouter()
	r.addRoute(http.MethodGet, "/a/*", mockHandler)
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	})

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
	})
}

// 這邊也是需要對 tree 的比較
// 返回一個錯誤訊息，幫助我們排查問題
// bool 是代表是否真的相等
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到對應的 http method: %s", k), false
		}
		// v 和 dst 要看是不是相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

// 實現一個 Equal 函數
func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("節點 path 不相等: %s != %s", n.path, y.path), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子節點數量不相等: %d != %d", len(n.children), len(y.children)), false
	}
	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, ok
		}
	}
	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}
	if n.regexChild != nil {
		msg, ok := n.regexChild.equal(y.regexChild)
		if !ok {
			return msg, ok
		}
	}

	// 比較 handler 是否相等
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return "handler 不相等", false
	}

	for path, child := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子節點 %s 不存在", path), false
		}
		msg, ok := child.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func TestRouter_FindRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		// {
		// 	method: http.MethodPost,
		// 	path:   "/login",
		// },
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/login/:username",
		},
		{
			method: http.MethodGet,
			path:   "/a/b/*",
		},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		info      *matchInfo
		// wantNode  *node
	}{
		{
			name:      "method not found",
			method:    http.MethodOptions,
			path:      "/",
			wantFound: false,
		},
		{
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "order star",
			method:    http.MethodGet,
			path:      "/order/abc",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "user home",
			method:    http.MethodGet,
			path:      "/user/home",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中但是沒有 handler
			name:      "order",
			method:    http.MethodGet,
			path:      "/order",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path: "order",
					children: map[string]*node{
						"detail": {
							path:    "detail",
							handler: mockHandler,
						},
					},
				},
			},
		},
		{
			// login 的參數路由
			name:      "login username",
			method:    http.MethodPost,
			path:      "/login/rex",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":username",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"username": "rex",
				},
			},
		},
		{
			// 根節點
			name:      "root",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
					children: map[string]*node{
						"order": {
							path: "order",
							children: map[string]*node{
								"detail": {
									path:    "detail",
									handler: mockHandler,
								},
							},
						},
						"user": {
							path: "user",
							children: map[string]*node{
								"home": {
									path:    "home",
									handler: mockHandler,
								},
							},
						},
						"a": {
							path: "a",
							children: map[string]*node{
								"b": {
									path: "b",
									starChild: &node{
										path:    "*",
										handler: mockHandler,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			// 测试 /a/b/* 匹配 /a/b/c
			name:      "a/b/c with wildcard",
			method:    http.MethodGet,
			path:      "/a/b/c",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "a/b/c/d with wildcard",
			method:    http.MethodGet,
			path:      "/a/b/c/d",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)

			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.info.pathParams, info.pathParams)
			msg, ok := tc.info.n.equal(info.n)
			assert.True(t, ok, msg)
		})
	}
}

func TestRouter_RegexParam(t *testing.T) {
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}

	// 测试正则表达式参数路由
	r.addRoute(http.MethodGet, "/user/:id(\\d+)", mockHandler)

	// 测试有效匹配
	info, found := r.findRoute(http.MethodGet, "/user/123")
	assert.True(t, found)
	assert.Equal(t, "123", info.pathParams["id"])

	// 测试无效匹配
	_, found = r.findRoute(http.MethodGet, "/user/abc")
	assert.False(t, found)

	// 测试冲突注册
	assert.Panic(t, func() {
		r.addRoute(http.MethodGet, "/user/:name([a-z]+)", mockHandler)
	})

	// 测试静态路由可以共存
	r.addRoute(http.MethodGet, "/user/profile", mockHandler)
	info, found = r.findRoute(http.MethodGet, "/user/profile")
	assert.True(t, found)
}

func TestParseParamRoute(t *testing.T) {
	// 测试正则表达式参数解析
	paramName, regexPattern, isParam := parseParamRoute(":id(\\d+)")
	fmt.Printf("paramName: %s, regexPattern: %s, isParam: %v\n", paramName, regexPattern, isParam)

	assert.True(t, isParam)
	assert.Equal(t, "id", paramName)
	assert.Equal(t, "\\d+", regexPattern)

	// 测试普通参数解析
	paramName, regexPattern, isParam = parseParamRoute(":id")
	fmt.Printf("paramName: %s, regexPattern: %s, isParam: %v\n", paramName, regexPattern, isParam)

	assert.True(t, isParam)
	assert.Equal(t, "id", paramName)
	assert.Equal(t, "", regexPattern)
}
