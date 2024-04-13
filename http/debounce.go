package http

import (
	"net/http"
	"sync"
	"time"
)

func DebounceFirst(circuit Circuit, duration time.Duration) Circuit {
	var (
		threshold time.Time
		response  *http.Response
		err       error
		mu        sync.Mutex
	)

	return func(r *http.Request) (*http.Response, error) {
		mu.Lock()

		defer func() {
			threshold = time.Now().Add(duration)
			mu.Unlock()
		}()

		if time.Now().Before(threshold) {
			return response, err
		}

		response, err = circuit(r)

		return response, err
	}
}
