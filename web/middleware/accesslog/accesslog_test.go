package accesslog

import (
	"github.com/rexshen5913/geek-learn-go/geektime-go /web"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	b := NewBuilder()
	s := web.NewHTTPServer()
	s.Get("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *web.Context) {
		time.Sleep(time.Second)
		ctx.RespData = []byte("hello, user")
	})
	s.UseAny("/*", b.Build())
	s.Start(":8081")
}
