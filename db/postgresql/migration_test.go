package postgresql

import (
	"context"
	"embed"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/igefined/go-kit/config"
	"github.com/igefined/go-kit/log"
	"github.com/igefined/go-kit/test"
)

func TestMakeMigrateUrl(t *testing.T) {
	tCases := []struct {
		name   string
		val    string
		result string
	}{
		{
			name:   "successfully, with sslmode",
			val:    "postgres://postgres:12345@localhost:5432/nh_templates?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			result: "postgres://postgres:12345@localhost:5432/nh_templates?sslmode=disable",
		},
		{
			name:   "successfully, with sslmode",
			val:    "postgres://postgres:12345@localhost:5432/nh_templates?pool_max_conns=16&sslmode=disable&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			result: "postgres://postgres:12345@localhost:5432/nh_templates?sslmode=disable",
		},
		{
			name:   "successfully, without sslmode",
			val:    "postgres://postgres:12345@localhost:5432/nh_templates?pool_max_conns=16&&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			result: "postgres://postgres:12345@localhost:5432/nh_templates?",
		},
	}

	for _, c := range tCases {
		t.Run(c.name, func(t *testing.T) {
			migrateUrl := makeMigrateUrl(c.val)
			assert.Equal(t, c.result, migrateUrl)
		})
	}
}

func TestMigrate(t *testing.T) {
	t.Run("invalid embed fs", func(t *testing.T) {
		logger, err := log.NewLogger(zap.DebugLevel)
		assert.NoError(t, err)

		cfg := &config.DBCfg{}

		err = Migrate(logger, &embed.FS{}, cfg)
		assert.ErrorContains(t, err, "file does not exist")
	})

	t.Run("error create source instance", func(t *testing.T) {
		logger, err := log.NewLogger(zap.DebugLevel)
		assert.NoError(t, err)

		cfg := &config.DBCfg{URL: "invalid_db_url"}

		err = Migrate(logger, &db, cfg)
		assert.ErrorContains(t, err, "no scheme")
	})

	t.Run("the schema not changed", func(t *testing.T) {
		logger, err := log.NewLogger(zap.DebugLevel)
		assert.NoError(t, err)

		var cfg = &config.DBCfg{
			URL: "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable",
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		container, err := test.NewPostgresContainer(ctx, cfg, &test.Opt{Enabled: true, Image: defaultPostgresImage})
		assert.NoError(t, err)
		assert.NotNil(t, container)

		err = Migrate(logger, &db, cfg)
		assert.NoError(t, err, err)

		err = Migrate(logger, &db, cfg)
		assert.ErrorIs(t, err, migrate.ErrNoChange)
	})

	t.Run("success", func(t *testing.T) {
		logger, err := log.NewLogger(zap.DebugLevel)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		var cfg = &config.DBCfg{
			URL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		}

		container, err := test.NewPostgresContainer(ctx, cfg, &test.Opt{Enabled: true, Image: defaultPostgresImage})
		assert.NoError(t, err)
		assert.NotNil(t, container)

		err = Migrate(logger, &db, cfg)
		assert.NoError(t, err)
	})
}
