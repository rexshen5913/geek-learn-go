package web

// Middleware 函數式的責任鏈模式
// 函數式的洋蔥模式
type Middleware func(next HandleFunc) HandleFunc

// AOP 切面編程，在不同的語言裡面都有不同的叫法
// Middlerware, Handler, Chain, Filter, Filter-Chain
// Interceptor, Wrapper
// type MiddlewareV1 interface {
// 	Invoke(next HandleFunc) HandleFunc
// }

// // 攔截器模式
// type Interceptor interface {
// 	Before(ctx *Context)
// 	After(ctx *Context)
// 	Surround(ctx *Context)
// }

// type HandlerFunV1 func(ctx *Context) (next bool)

// type Chain []HandlerFunV1

// type ChainV1 struct {
// 	handlers []HandlerFunV1
// }

// func (c ChainV1) Run(ctx *Context) {
// 	for _, handler := range c.handlers {
// 		next := handler(ctx)
// 		// 這種是中斷執行
// 		if !next {
// 			return
// 		}
// 	}
// }

// type Net struct {
// 	handlers []HandleFuncV1
// }

// func (c Net) Run(ctx *Context) {
// 	wg := sync.WaitGroup{}
// 	for _, handler := range c.handlers {
// 		h := handler
// 		if h.concurrent {
// 			wg.Add(1)
// 			go func() {
// 				h.Run(ctx)
// 				wg.Done()
// 			}()
// 		} else {
// 			h.Run(ctx)
// 		}
// 	}
// 	wg.Wait()
// }

// type HandleFuncV1 struct {
// 	concurrent bool
// 	handlers   []HandleFuncV1
// }

// func (h HandleFuncV1) Run(ctx *Context) {
// 	// for _, handler := range h.handlers {
// 	// 	handler.Run(ctx)
// 	// }
// }
