package accesslog

import (
	"fmt"
	"log"
	"testing"
	"web"
)

func TestMiddlewareBuilder_E2E(t *testing.T) {
	// 鏈式調用
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()

	// 這裡是使用 ServerWithMiddleware 來追加 middleware
	server := web.NewHttpServer(web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, It's me"))
	})

	log.Println("server start at :8081")
	server.Start(":8081")

}
