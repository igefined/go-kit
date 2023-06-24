//go:build units

package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igdotog/core/config"
	baseHttp "github.com/igdotog/core/http"
	"github.com/igdotog/core/logger"

	asserting "github.com/stretchr/testify/assert"
)

func TestNewBaseServer(t *testing.T) {
	assert := asserting.New(t)
	cfg := config.Application{
		Http: &config.Http{
			Port: 9977,
		},
	}

	logger := logger.New()

	t.Run("construct without fx", func(t *testing.T) {
		bs := baseHttp.NewBaseServer(&cfg, logger, nil)
		assert.NotNil(bs)
	})

	t.Run("check health", func(t *testing.T) {
		const expectText = "ok"

		req := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
		rec := httptest.NewRecorder()

		baseHttp.HealthHandler().ServeHTTP(rec, req)

		res := rec.Result()
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				t.Log(err)
			}
		}(res.Body)

		resBytes, err := io.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(string(resBytes), expectText)
		assert.Equal(res.StatusCode, http.StatusOK)
	})
}
