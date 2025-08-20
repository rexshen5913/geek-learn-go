package web

import (
	"fmt"
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
	addRoute(method string, path string, handleFunc HandleFunc, mdls ...Middleware)
	// AddRoute1(method string, path string, handleFunc ...HandleFunc)
}

// type HTTPServer struct {
// 	HttpServer
// }

type HTTPServerOption func(server *HttpServer)

type HttpServer struct {
	router
	middlewares []Middleware

	log func(msg string, args ...any)
}

func NewHttpServer(opts ...HTTPServerOption) *HttpServer {
	server := &HttpServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

func ServerWithMiddleware(middlewares ...Middleware) HTTPServerOption {
	return func(server *HttpServer) {
		server.middlewares = middlewares
	}
}

// ServeHTTP 实现 Server 接口
func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 你的框架代碼寫在這裡
	ctx := &Context{
		Resp: writer,
		Req:  request,
	}

	// 最後一個是這個
	root := h.serve
	// 利用最後一個不斷往前回朔組裝鏈條
	// 從後往前
	// 把後一個作為前一個的 next
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		root = h.middlewares[i](root)
	}

	// 這裡就是執行的時候，就是從前往後了

	// 這裡最後一個步驟，就是把 RespData 和 RespStatusCode 刷新到響應中
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			// 就設置好了 RespData 和 RespStatusCode
			next(ctx)
			h.flushResp(ctx)
		}
	}

	// 因為會需要最後一個執行
	root = m(root)

	root(ctx)
	// h.serve(ctx)
}

func (h *HttpServer) flushResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		// log.Fatalln("web: 寫入響應數據失敗", err)
		h.log("web: 寫入響應數據失敗 %v", err)
	}
}

func (h *HttpServer) serve(ctx *Context) {
	// 接下來就是查找路由並執行命中的業務邏輯
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		// ctx.Resp.WriteHeader(http.StatusNotFound)
		// _, _ = ctx.Resp.Write([]byte("Not Found"))
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("Not Found")
		return
	}

	ctx.PathParams = info.pathParams
	// 應該是要叫 pattern 比較準確，但一般都是叫 path
	ctx.MatchedRoute = info.n.route

	root := info.n.handler

	for i := len(info.mdls) - 1; i >= 0; i-- {
		root = info.mdls[i](root)
	}

	root(ctx)
}

func (h *HttpServer) Use(method string, path string, mdls ...Middleware) {
	h.addRoute(method, path, nil, mdls...)

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
