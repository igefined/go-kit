//go:build units

package config_test

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/igdotog/core/config"

	testAssert "github.com/stretchr/testify/assert"
)

const (
	appName      = "application"
	defaultDBUrl = "postgres://postgres:12345@localhost:5432/testdb?sslmode=disable"
)

func TestNew(t *testing.T) {
	assert := testAssert.New(t)
	ctx := context.Background()

	setEnvHelper(t, "DEBUG", "true")
	setEnvHelper(t, "APP_NAME", appName)

	t.Run("base config", func(t *testing.T) {
		cfg, err := config.New(ctx)
		assert.NoError(err)
		assert.Equal(cfg.Name, appName)
		assert.True(strings.HasPrefix(cfg.AppId, appName))
		assert.Nil(cfg.DB)
	})

	t.Run("db config", func(t *testing.T) {
		setEnvHelper(t, "DB_URL", defaultDBUrl)
		setEnvHelper(t, "DB_AUTO_CREATE_DATABASE", "true")

		cfg, err := config.New(ctx, &config.DB{})
		assert.NoError(err)
		assert.Equal(cfg.Name, appName)
		assert.True(strings.HasPrefix(cfg.AppId, appName))
		assert.Equal(cfg.DB.URL, defaultDBUrl)
		assert.True(cfg.DB.AutoCreateDatabase)
	})

	t.Run("http config", func(t *testing.T) {
		httpPort := 3000
		setEnvHelper(t, "HTTP_PORT", strconv.Itoa(httpPort))

		cfg, err := config.New(ctx, &config.Http{})
		assert.NoError(err)
		assert.Equal(cfg.Name, appName)
		assert.True(strings.HasPrefix(cfg.AppId, appName))
		assert.Equal(cfg.Http.Port, uint(httpPort))
	})
}

func setEnvHelper(t *testing.T, key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Fatal(err)
	}
}
