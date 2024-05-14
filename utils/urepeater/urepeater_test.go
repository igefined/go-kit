//go:build units

package urepeater

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepeaterSuccess(t *testing.T) {
	callCount := 0
	f := func() (string, error) {
		callCount++
		return "success", nil
	}

	res, err := Repeater(f)
	assert.NoError(t, err)
	assert.Equal(t, "success", res)
	assert.Equal(t, 1, callCount, "Function should be called only once")
}

func TestRepeaterRetries(t *testing.T) {
	callCount := 0
	f := func() (string, error) {
		callCount++
		if callCount < 2 {
			return "", errors.New("temporary error")
		}
		return "success", nil
	}

	sleep = func(d time.Duration) {
		// Do nothing
	}

	res, err := Repeater(f)
	assert.NoError(t, err)
	assert.Equal(t, "success", res)
	assert.Equal(t, 2, callCount, "Function should be called two times")
}

var sleep = time.Sleep

func init() {
	rand.New(rand.NewSource(1))
}
