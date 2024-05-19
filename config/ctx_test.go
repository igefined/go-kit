package config

import (
	"os"
	"syscall"
	"testing"
	"time"

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

	t.Run("success with term signal", func(t *testing.T) {
		initCtx := SigTermIntCtx()
		assert.NotNil(t, initCtx)
	})
}

func TestSigTermIntCtx_NoSignal(t *testing.T) {
	tCtx := SigTermIntCtx()

	select {
	case <-tCtx.Done():
		t.Fatal("context should not be cancelled")
	default:
	}
}

func TestSigTermIntCtx_SIGTERM(t *testing.T) {
	tCtx := SigTermIntCtx()

	pid := os.Getpid()

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(pid, syscall.SIGTERM)
	}()

	select {
	case <-tCtx.Done():
	case <-time.After(1 * time.Second):
		t.Fatal("context was not cancelled after SIGTERM")
	}
}

func TestSigTermIntCtx_SIGINT(t *testing.T) {
	tCtx := SigTermIntCtx()

	pid := os.Getpid()

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(pid, syscall.SIGINT)
	}()

	select {
	case <-tCtx.Done():
	case <-time.After(1 * time.Second):
		t.Fatal("context was not cancelled after SIGINT")
	}
}
