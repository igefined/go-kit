package uthrottler

import (
	"context"
	"time"
)

type Effector func(ctx context.Context) (string, error)

type Throttled func(context.Context, string) (bool, string, error)

type bucket struct {
	tokens uint
	time   time.Time
}

func Throttle(e Effector, max, refill uint, d time.Duration) Throttled {
	buckets := map[string]*bucket{}

	return func(ctx context.Context, key string) (bool, string, error) {
		b, ok := buckets[key]
		if !ok {
			buckets[key] = &bucket{tokens: max - 1, time: time.Now()}

			str, err := e(ctx)
			return true, str, err
		}

		refillInterval := uint(time.Since(b.time) / d)
		tokensAdded := refill * refillInterval
		currentTokens := b.tokens + tokensAdded

		if currentTokens < 1 {
			return false, "", nil
		}

		if currentTokens > max {
			b.time = time.Now()
			b.tokens = max - 1
		} else {
			deltaTokens := currentTokens - b.tokens
			deltaRefills := deltaTokens / refill
			deltaTime := time.Duration(deltaRefills) * d

			b.time = b.time.Add(deltaTime)
			b.tokens = currentTokens - 1
		}

		str, err := e(ctx)

		return true, str, err
	}
}
