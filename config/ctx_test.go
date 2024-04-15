package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigTermIntCtx(t *testing.T) {
	t.Run("success sig term context", func(t *testing.T) {
		initCtx := SigTermIntCtx()
		assert.NotNil(t, initCtx)

		doubleCtx := SigTermIntCtx()
		assert.NotNil(t, doubleCtx)
		assert.Equal(t, initCtx, doubleCtx)
	})
}
