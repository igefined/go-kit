package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// ConstructNewSpanContext helper function to inject in open telemetry refid from context
func ConstructNewSpanContext(ctx context.Context) trace.SpanContext {
	var (
		spanContext trace.SpanContext
		spanID      trace.SpanID
		traceID     trace.TraceID
	)

	traceIDStr, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		traceIDStr = GenerateTraceIDString()
	}

	traceID, err := trace.TraceIDFromHex(traceIDStr)
	if err != nil {
		traceID, _ = trace.TraceIDFromHex(GenerateTraceIDString())
	}

	spanIDStr, ok := ctx.Value(SpanIDKey).(string)
	if !ok {
		spanIDStr = GenerateSpanIDString()
	}

	spanID, err = trace.SpanIDFromHex(spanIDStr)
	if err != nil {
		spanID, _ = trace.SpanIDFromHex(GenerateSpanIDString())
	}

	spanContext = trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: 0x1,
	})
	return spanContext
}

// HTTPMiddleware Assume that spanID and TraceID already in ctx
// Example:
//
//	h := refid.HTTPMiddleware(tracing.HTTPMiddleware(otelhttp.NewHandler(s.newRouter(), "app-test",
//		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
//		otelhttp.WithPropagators(propagation.TraceContext{}),
//	)))
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		xTraceID, ok := ctx.Value(TraceIDKey).(string)
		if ok {
			w.Header().Add(TraceIDKey, xTraceID)
		}

		next.ServeHTTP(w, r.WithContext(trace.ContextWithSpanContext(r.Context(), ConstructNewSpanContext(r.Context()))))
	})
}
