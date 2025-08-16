package opentelemetry

import (
	"web"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/rexshen5913/geek-learn-go/geektime-go/Mytest/web/middleware/opentelemetry"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {

			// 嘗試和客戶端的 trace 建立關聯
			reqCtx := ctx.Req.Context()
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			// 創建一個span
			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")

			defer func() {
				// 這是只有執行完 next 才可能有值
				if ctx.MatchedRoute != "" {
					span.SetName(ctx.MatchedRoute)
				}

				// 把響應碼加上去
				span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
				span.End()
			}()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.schema", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))

			// 將 trace 的 context 傳給下一個 middleware
			ctx.Req = ctx.Req.WithContext(reqCtx)

			// 二話不說，先直接調用下一個middleware
			// 確保忘記
			next(ctx)
		}
	}
}
