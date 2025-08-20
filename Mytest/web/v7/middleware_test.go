package web

import (
	"testing"
)

func TestMiddlewareAppend(t *testing.T) {
	// 測試追加 middleware 的行為
	server := NewHttpServer(
		ServerWithMiddleware(func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				// middleware 1
				next(ctx)
			}
		}),
		ServerWithMiddleware(func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				// middleware 2
				next(ctx)
			}
		}),
	)

	// 應該有 2 個 middleware
	if len(server.middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(server.middlewares))
	}
}

func TestMiddlewareOverride(t *testing.T) {
	// 測試覆蓋 middleware 的行為
	server := NewHttpServer(
		ServerWithMiddleware(func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				// 這個會被覆蓋
				next(ctx)
			}
		}),
		ServerWithMiddleware(func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				// 這個會保留
				next(ctx)
			}
		}),
	)

	// 應該只有 1 個 middleware（被覆蓋了）
	if len(server.middlewares) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(server.middlewares))
	}
}

func TestMiddlewareAppendMultiple(t *testing.T) {
	// 測試多次追加多個 middleware
	server := NewHttpServer(
		ServerWithMiddleware(
			func(next HandleFunc) HandleFunc {
				return func(ctx *Context) { next(ctx) }
			},
			func(next HandleFunc) HandleFunc {
				return func(ctx *Context) { next(ctx) }
			},
		),
		ServerWithMiddleware(
			func(next HandleFunc) HandleFunc {
				return func(ctx *Context) { next(ctx) }
			},
		),
	)

	// 應該有 3 個 middleware
	if len(server.middlewares) != 3 {
		t.Errorf("Expected 3 middlewares, got %d", len(server.middlewares))
	}
}
