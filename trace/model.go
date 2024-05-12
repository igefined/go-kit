package trace

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceIDKey = "x-trace-id"
	SpanIDKey  = "x-span-id"
)

type TraceID struct {
}

func GenerateTraceIDRaw() []byte {
	id, _ := uuid.New().MarshalBinary()
	return id
}

func GenerateTraceIDString() string {
	return hex.EncodeToString(GenerateTraceIDRaw())
}

func (t *TraceID) Generate() any {
	return GenerateTraceIDString()
}

func (t *TraceID) IsValid(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := trace.TraceIDFromHex(s)
	return err == nil
}

type SpanID struct {
}

func GenerateSpanIDRaw() []byte {
	id := make([]byte, 8)
	_, _ = rand.Read(id)
	return id
}

func GenerateSpanIDString() string {
	return hex.EncodeToString(GenerateSpanIDRaw())
}

func (t *SpanID) Generate() any {
	return GenerateSpanIDString()
}

func (t *SpanID) IsValid(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := trace.SpanIDFromHex(s)
	return err == nil
}
