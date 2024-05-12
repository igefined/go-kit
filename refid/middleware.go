package refid

import (
	"context"
	"net/http"
	"strings"
)

const (
	tracePrefix = "x-trace-"
	spanPrefix  = "x-span-"
)

func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := ctxTags{}
		for key, values := range r.Header {
			if len(values) == 0 {
				continue
			}

			lKey := strings.ToLower(key)
			_, ok := markers[lKey]
			if ok || strings.HasPrefix(lKey, tracePrefix) || strings.HasPrefix(lKey, spanPrefix) {
				tags[lKey] = values[0]
			}
		}

		for k, v := range markers {
			if _, ok := tags[k]; !ok {
				tags[k] = v.Generate()
			}
		}

		next.ServeHTTP(w, r.WithContext(PutTags(r.Context(), tags)))
	})
}

func AddToHTTPClientHeaders(ctx context.Context, override map[string]string, header *http.Header) {
	for k, v := range GetFields(ctx) {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if val, ok := override[k]; ok {
			header.Add(k, val)
			continue
		}
		_, ok = markers[k]
		if ok || strings.HasPrefix(k, tracePrefix) || strings.HasPrefix(k, spanPrefix) {
			header.Add(k, s)
		}
	}
}
