package accesslog

import (
	"fmt"
	"net/http"
	"testing"
	"web"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	// 鏈式調用
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()

	// 這裡是使用 ServerWithMiddleware 來追加 middleware
	server := web.NewHttpServer(web.ServerWithMiddleware(mdl))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, It's me")
	})

	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
