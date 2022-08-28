package iris

import (
	"github.com/kataras/iris/v12"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.HTML("Hello <strong>%s</strong>!", "World")
	})

	_ = app.Listen(":8083")
}
