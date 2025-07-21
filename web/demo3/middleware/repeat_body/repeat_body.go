package repeat_body

import (
	web "github.com/rexshen5913/geek-learn-go/geektime-go /web/demo3"
	"io/ioutil"
)

func Middleware() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			ctx.Req.Body = ioutil.NopCloser(ctx.Req.Body)
			next(ctx)
		}
	}
}
