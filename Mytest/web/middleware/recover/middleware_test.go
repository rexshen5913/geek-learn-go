package recover

import (
	"fmt"
	"net/http"
	"testing"
	"web"
)

func TestMiddlewareBuild_Build(t *testing.T) {
	builder := &MiddlewareBuilder{
		StatusCode: http.StatusInternalServerError,
		Data:       []byte("Got panic !!!"),
		Log: func(ctx *web.Context) {
			fmt.Printf("panic 路徑: %s\n", ctx.Req.URL.String())
		},
	}
	server := web.NewHttpServer(web.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *web.Context) {
		panic("user not found")
	})
	server.Start(":8081")
}
