package web

import (
	"fmt"
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

	h.addRoute("GET", "/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello, order detail"))
	})

	handler1 := func(ctx *Context) {
		fmt.Println("處理第一件事!")
	}

	handler2 := func(ctx *Context) {
		fmt.Println("處理第二件事!")
	}

	h.Get("/user1", handler1)
	h.Get("/user2", handler2)

	// 調試點：在啟動服務器前暫停
	fmt.Println("=== 調試點 2: 準備啟動服務器 ===")
	fmt.Println("按 Enter 繼續...")

	// 驗證 h 實現了 Server 接口
	var _ Server = h
	fmt.Println("server start at 8081")
	h.Start(":8081")
}
