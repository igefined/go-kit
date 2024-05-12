package trace

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

// NewGRPCClient create new open telemetry tracing grpc client
func NewGRPCClient(collectorURI string, options ...otlptracegrpc.Option) otlptrace.Client { //nolint:ireturn
	options = append(
		options,
		otlptracegrpc.WithEndpoint(collectorURI),
	)

	if len(options) == 0 {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.NewClient(options...)
}

// NewHTTPClient create new open telemetry tracing http client
func NewHTTPClient(collectorURI string, options ...otlptracehttp.Option) otlptrace.Client {
	options = append(
		options,
		otlptracehttp.WithEndpoint(collectorURI),
	)

	if len(options) == 0 {
		options = append(options, otlptracehttp.WithInsecure())
	}

	return otlptracehttp.NewClient(options...)
}
