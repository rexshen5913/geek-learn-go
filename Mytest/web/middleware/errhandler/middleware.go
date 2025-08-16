package errhandler

import "web"

type MiddlewareBuilder struct {
	// 這種設計只能返回固定的值
	// 不能做到動態渲染
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: make(map[int][]byte, 64),
	}
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				resp, ok := m.resp[ctx.RespStatusCode]
				if ok {
					// 竄改結果
					ctx.RespData = resp
					// 设置正确的 Content-Type
					ctx.Resp.Header().Set("Content-Type", "text/html")
				}

			}()
			next(ctx)
		}
	}
}
