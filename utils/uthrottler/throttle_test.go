//go:build units

package uthrottler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type uKeyCtx struct{}

var errUniqueKeyRequired = errors.New("unique key required")

func TestThrottle(t *testing.T) {
	const testIP = "127.0.0.1"

	ctx := context.WithValue(context.Background(), uKeyCtx{}, testIP)

	effector := func(ctx context.Context) (string, error) {
		val, ok := ctx.Value(uKeyCtx{}).(string)
		if !ok {
			return "", errUniqueKeyRequired
		}

		return val, nil
	}

	throttle := Throttle(effector, 1, 1, time.Second)

	b, s, err := throttle(ctx, testIP)
	assert.NoError(t, err)
	assert.True(t, b)
	assert.NotEmpty(t, s)

	b, s, err = throttle(ctx, testIP)
	assert.NoError(t, err, errUniqueKeyRequired)
	assert.False(t, b)
	assert.Empty(t, s)
}
