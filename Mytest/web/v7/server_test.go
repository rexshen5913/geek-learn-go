package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHttpServer()
	server.middlewares = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一個 before")
				next(ctx)
				fmt.Println("第一個 after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二個 before")
				next(ctx)
				fmt.Println("第二個 after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三個中斷")
				// next(ctx)
				// fmt.Println("第三個 after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四個你看不到這句話")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
