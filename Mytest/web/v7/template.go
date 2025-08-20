package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	// Render 渲染模板
	// tplName 模板名稱, 按名稱索引
	// data 是渲染模板所需數據
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	// 渲染模板並寫入 writer
	// Render(ctx context.Context, tplName string, data any, writ er io.Writer) error

	// 不需要，讓具體的實現自己管自己的模板
	// AddTemplate(name string, tplContent []byte) error
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}
