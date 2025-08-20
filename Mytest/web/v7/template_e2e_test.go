package web

import (
	"html/template"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoginPage(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.html")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}
	h := NewHttpServer(ServerWithTemplateEngine(engine))

	h.Get("/login", func(ctx *Context) {
		err := ctx.Render("login.html", map[string]any{
			"Username": "Rex",
			"Email":    "rex@example.com",
		})
		if err != nil {
			log.Println(err)
		}
	})

	h.Start(":8080")
}
