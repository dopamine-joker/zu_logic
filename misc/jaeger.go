package misc

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strings"
)

const (
	SERVICE     = "zu-logic"
	environment = "production"
)

var Tracer trace.Tracer

// tracerProvide new a tracerProvider
func tracerProvide(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(SERVICE),
			attribute.String("environment", environment),
		)),
	)
	return tp, nil
}

//initJaeger jaeger初始化
func initJaeger() {
	provider, err := tracerProvide(strings.Join([]string{Conf.JaegerCfg.Schema, Conf.JaegerCfg.Host,
		Conf.JaegerCfg.Path}, ""))
	if err != nil {
		Logger.Error("init tracerProvider err", zap.String("err", err.Error()))
		panic(err)
	}
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	Tracer = otel.Tracer("zu-logic")
}
