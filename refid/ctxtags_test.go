package refid

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	customKey    = "x-trace-custom"
	customValue1 = "custom_value"

	traceIDKey = "x-trace-id"
	spanIDKey  = "x-span-id"
)

func TestSetGet(t *testing.T) {
	ctx := PutTags(context.Background(), map[string]any{
		"key1": "val1",
	})

	ctx = PutTag(ctx, "key2", "val2")
	ctx = PutTag(ctx, "key3", "val3")
	ctx = PutTag(ctx, "key1", "over")
	val2, ok := GetStringField(ctx, "key2")
	require.True(t, ok)
	val3, ok := GetStringField(ctx, "key3")
	require.True(t, ok)
	val1, ok := GetStringField(ctx, "key1")
	require.True(t, ok)

	require.Equal(t, "val2", val2)
	require.Equal(t, "val3", val3)
	require.Equal(t, "over", val1)
}

func TestCheckHTTPMiddlewareExisting(t *testing.T) {
	traceID, _ := uuid.New().MarshalBinary()
	traceIDStr := hex.EncodeToString(traceID)
	spanID := "1234"

	check := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s1, ok := GetStringField(r.Context(), traceIDKey)
			require.True(t, ok)
			require.Equal(t, traceIDStr, s1)
			span, ok := GetStringField(r.Context(), spanIDKey)
			require.True(t, ok)
			require.Equal(t, spanID, span)

			customValueGet, ok := GetStringField(r.Context(), customKey)
			require.True(t, ok)
			require.Equal(t, customValue1, customValueGet)
		}
	}

	ctx := PutTags(context.Background(), ctxTags{
		traceIDKey: traceIDStr,
		spanIDKey:  spanID,
		customKey:  customValue1,
	})

	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	AddToHTTPClientHeaders(ctx, nil, &req.Header)
	rr := httptest.NewRecorder()
	HTTPMiddleware(check()).ServeHTTP(rr, req)
}
