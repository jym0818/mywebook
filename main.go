package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {

	initViper()
	initPrometheus()
	fn := initOpentelemetry()

	app := InitWeb()

	cron := app.cron
	cron.Start()

	server := app.web
	server.Run(":8080")
	//10s内必须关闭否则强制关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fn(ctx)
	//可以考虑强制退出，防止有些任务执行时间太长
	stop := cron.Stop()
	timer := time.NewTimer(5 * time.Minute)
	select {
	case <-timer.C:
	case <-stop.Done():
	}
}

func initViper() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func initOpentelemetry() func(ctx context.Context) {
	res, err := newResource("demo", "v0.0.1")
	if err != nil {
		panic(err)
	}

	prop := newPropagator()
	// 在客户端和服务端之间传递 tracing 的相关信息
	otel.SetTextMapPropagator(prop)

	// 初始化 trace provider
	// 这个 provider 就是用来在打点的时候构建 trace 的
	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}
	//	defer tp.Shutdown(context.Background())
	//这里如何close
	otel.SetTracerProvider(tp)
	return func(ctx context.Context) {
		// 这个ctx是用于超时控制，如果超过时间还没有关闭 会强制关闭
		tp.Shutdown(ctx)
	}
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
