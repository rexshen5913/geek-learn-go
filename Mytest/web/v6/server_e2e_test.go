package web

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestServer(t *testing.T) {
	// 調試點：程序會在這裡暫停
	fmt.Println("=== 調試點 1: 開始測試 ===")

	h := NewHttpServer()

	h.addRoute("GET", "/user", func(ctx *Context) {
		fmt.Println("處理第一件事!")
		fmt.Println("處理第二件事!")
	})

	h.addRoute("GET", "/order/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("經由通配符路由"))
	})

	h.addRoute("GET", "/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello, order detail"))
	})

	// 測試通配符路由
	h.addRoute("GET", "/a/b/*", func(ctx *Context) {
		// 這邊會回傳完整路由路徑
		ctx.Resp.Write([]byte("經由通配符路由: " + ctx.Req.URL.Path))
	})

	// 測試正則表達式參數路由
	h.addRoute("GET", "/user/:id(\\d+)", func(ctx *Context) {
		ctx.Resp.Write([]byte("經由正則表達式參數路由: " + ctx.Req.URL.Path))
	})

	handler1 := func(ctx *Context) {
		fmt.Println("處理第一件事!")
	}

	handler2 := func(ctx *Context) {
		fmt.Println("處理第二件事!")
	}

	h.Get("/user1", handler1)
	h.Get("/user2", handler2)

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	h.Get("/user/123", func(ctx *Context) {
		ctx.RespStatusOK(User{
			Name: "John",
			Age:  18,
		})
	})

	h.Get("/user/456", func(ctx *Context) {
		safeCtx := SafeContext{
			ctx: ctx,
		}
		safeCtx.RespStatusOK(User{
			Name: "Rex",
			Age:  28,
		})
	})

	h.Get("/user/info", func(ctx *Context) {
		age, err := QueryValueV3[int](ctx, "age", strconv.Atoi)
		if err != nil {
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			ctx.Resp.Write([]byte("bad request"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %d", age)))
	})

	h.Post("/form", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, " + ctx.Req.URL.Path))
	})

	h.Get("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			ctx.Resp.Write([]byte("bad request"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %d", id)))
	})
	// 調試點：在啟動服務器前暫停
	fmt.Println("=== 調試點 2: 準備啟動服務器 ===")
	fmt.Println("按 Enter 繼續...")

	// 驗證 h 實現了 Server 接口
	var _ Server = h
	fmt.Println("server start at 8081")
	h.Start(":8081")
}
