package opentelemetry

import (
	"context"
	"fmt"
	"testing"
	"time"
	"web"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initTracer() (func(), error) {
	// 連接到 OpenTelemetry Collector
	conn, err := grpc.Dial("192.168.1.54:14317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	// 創建 OTLP exporter
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 創建 resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-web-server"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 創建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 設置全局 TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 返回清理函數
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

func Test_NewMiddlewareBuilder(t *testing.T) {
	// 初始化 tracer
	cleanup, err := initTracer()
	if err != nil {
		t.Fatalf("Failed to init tracer: %v", err)
	}
	defer cleanup()

	// 創建 middleware builder
	builder := &MiddlewareBuilder{
		Tracer: otel.GetTracerProvider().Tracer(instrumentationName),
	}

	// 創建服務器
	server := web.NewHttpServer(web.ServerWithMiddleware(builder.Build()))

	// 註冊路由
	server.Get("/user", func(ctx *web.Context) {
		// 在 handler 中創建一個子 span
		c, firstSpan1 := builder.Tracer.Start(ctx.Req.Context(), "first_layer_1")
		defer firstSpan1.End()

		secondC, secondSpan := builder.Tracer.Start(c, "second_layer_1")
		time.Sleep(100 * time.Millisecond)

		_, thirdSpan1 := builder.Tracer.Start(secondC, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		_, thirdSpan2 := builder.Tracer.Start(secondC, "third_layer_2")
		time.Sleep(100 * time.Millisecond)
		thirdSpan2.End()
		thirdSpan1.End()
		secondSpan.End()

		_, firstSpan2 := builder.Tracer.Start(ctx.Req.Context(), "first_layer_2")
		time.Sleep(100 * time.Millisecond)
		defer firstSpan2.End()

		time.Sleep(100 * time.Millisecond) // 模擬一些處理時間
		ctx.RespJSON(202, User{
			Name: "John Doe",
		})

	})
	server.Start(":8081")
}

type User struct {
	Name string
}
