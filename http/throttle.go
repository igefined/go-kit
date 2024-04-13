package http

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

var ErrTooManyRequests = errors.New("too many requests")

type Effector func(r *http.Request) (*http.Response, error)

func Throttle(e Effector, max, refill uint, d time.Duration) Effector {
	var (
		tokens = max
		once   sync.Once
	)

	return func(r *http.Request) (*http.Response, error) {
		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-r.Context().Done():
						return
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})

		if tokens <= 0 {
			return nil, ErrTooManyRequests
		}

		tokens--

		return e(r)
	}
}
