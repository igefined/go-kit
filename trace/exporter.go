package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewTraceSpanExporter(ctx context.Context, client otlptrace.Client) (tracesdk.SpanExporter, error) {
	return otlptrace.New(ctx, client)
}
