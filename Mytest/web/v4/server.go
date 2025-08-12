package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 驗證 HttpServer 實現了 Server 接口
var _ Server = &HttpServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	// 增加路由註冊的功能
	// method 是 http 方法，path 是路徑，handlerFunc 是處理函數
	addRoute(method string, path string, handleFunc HandleFunc)
	// AddRoute1(method string, path string, handleFunc ...HandleFunc)
}

// type HTTPServer struct {
// 	HttpServer
// }

type HttpServer struct {
	router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		router: newRouter(),
	}
}

// ServeHTTP 实现 Server 接口
func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 你的框架代碼寫在這裡
	ctx := &Context{
		Resp: writer,
		Req:  request,
	}

	h.serve(ctx)
}

func (h *HttpServer) serve(ctx *Context) {
	// 接下來就是查找路由並執行命中的業務邏輯
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("Not Found"))
		return
	}

	ctx.PathParams = info.pathParams

	info.n.handler(ctx)
}

// func (h *HttpServer) addRoute(method string, path string, handleFunc HandleFunc) {
// 	// 這裡可以做路由的註冊
// 	// panic("implemented me")
// }

func (h *HttpServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *HttpServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

func (h *HttpServer) Options(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}

// func (h *HttpServer) Put(path string, handleFunc HandleFunc) {
// 	h.addRoute(http.MethodPut, path, handleFunc)
// }

// func (h *HttpServer) Delete(path string, handleFunc HandleFunc) {
// 	h.addRoute(http.MethodDelete, path, handleFunc)
// }

// func (h *HttpServer) AddRoute1(method string, path string, handleFunc ...HandleFunc) {
// 	panic("implemented me")
// }

func (h *HttpServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 在這裡，可以讓用戶註冊所謂的 after start 回調
	// 比如說往你的 admin 註冊一下自己的這個實例
	// 在這裡執行一些你的業務所需的前置條件
	return http.Serve(l, h)
}

func (h *HttpServer) Start1(addr string) error {
	return http.ListenAndServe(addr, h)
}
