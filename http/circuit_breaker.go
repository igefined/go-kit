package http

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

var ErrCircuitBreaker = errors.New("service unreachable")

type Circuit func(r *http.Request) (*http.Response, error)

func Breaker(circuit Circuit, failureThreshold int) Circuit {
	var (
		consecutiveFailures = 0
		lastAttempt         = time.Now()
		mu                  sync.RWMutex
	)

	return func(r *http.Request) (*http.Response, error) {
		mu.RLock()

		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				mu.RUnlock()

				return nil, ErrCircuitBreaker
			}
		}

		mu.RUnlock()

		response, err := circuit(r)

		mu.Lock()
		defer mu.Unlock()
		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailures++
			return nil, err
		}

		consecutiveFailures = 0

		return response, nil
	}
}
