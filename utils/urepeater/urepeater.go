package urepeater

import (
	"math/rand"
	"time"
)

func Repeater[T any](f func() (T, error)) (T, error) {
	res, err := f()
	base, capacity := time.Second, time.Minute

	for backoff := base; err != nil; backoff <<= 1 {
		if backoff > base {
			backoff = capacity
		}

		jitter := rand.Int63n(int64(backoff * 3))
		sleep := base + time.Duration(jitter)
		time.Sleep(sleep)
		res, err = f()
	}

	return res, err
}
