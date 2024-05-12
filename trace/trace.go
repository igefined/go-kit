package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewGlobalTracerProvider(
	serviceName string,
	spanExporter sdk.SpanExporter,
	providerOpts ...sdk.TracerProviderOption,
) {
	opts := []sdk.TracerProviderOption{
		sdk.WithBatcher(spanExporter),
		sdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	}

	opts = append(opts, providerOpts...)
	tp := sdk.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
