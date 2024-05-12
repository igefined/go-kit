package postgresql

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/igefined/go-kit/config"
	"github.com/igefined/go-kit/log"
	"github.com/igefined/go-kit/test"
)

const (
	defaultPostgresImage = "postgres:15.3-alpine"

	prefix = "test"
)

type Suite struct {
	suite.Suite
	ctx context.Context

	cfg       *config.DBCfg
	logger    *log.Logger
	container *test.PostgresContainer
}

func TestBuilderSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	logger, err := log.NewLogger(zap.DebugLevel)
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	_ = cancel

	type testConfig struct {
		sync.RWMutex
		config.DBCfg `mapstructure:",squash"`
	}

	var cfg *testConfig
	s.Require().NoError(config.GetConfig(prefix, &cfg, []*config.EnvVar{}))
	s.Require().NotEmpty(cfg.URL)

	s.logger = logger
	s.cfg = &cfg.DBCfg
	s.ctx = ctx

	pgContainer, err := test.NewPostgresContainer(ctx, s.cfg, &test.Opt{Enabled: true, Image: defaultPostgresImage})
	s.Require().NoError(err)

	s.container = pgContainer
}

func (s *Suite) TearDownSuite() {
	if err := s.container.Terminate(s.ctx); err != nil {
		s.logger.Error("error terminating postgres container", zap.Error(err))
	}
}

func (s *Suite) TestCreateAndDropDatabase() {
	var isExists bool

	CreateDatabase(s.ctx, s.logger, s.cfg)

	pool, err := pgxpool.New(s.ctx, ReplaceDbName(s.cfg.URL, "postgres"))
	s.Require().NoError(err)
	defer pool.Close()

	row := pool.QueryRow(s.ctx, checkingSql, s.cfg.GetDatabaseName())
	err = row.Scan(&isExists)
	s.Require().NoError(err)
	s.Require().True(isExists)

	dropDbSql := fmt.Sprintf("drop database if exists %s", s.cfg.GetDatabaseName())
	rows, err := pool.Query(s.ctx, dropDbSql)
	s.Require().NoError(err)
	s.Require().NotNil(rows)
	rows.Close()
}

func TestReplaceDbName(t *testing.T) {
	tCases := []struct {
		srcUrl    string
		dbName    string
		resultUrl string
	}{
		{
			srcUrl:    "postgres://postgres:postgres@localhost:5466/test?sslmode=disable",
			dbName:    "silly",
			resultUrl: "postgres://postgres:postgres@localhost:5466/silly?sslmode=disable",
		},
		{
			srcUrl:    "postgres://postgres:postgres@localhost:5466/test",
			dbName:    "silly",
			resultUrl: "postgres://postgres:postgres@localhost:5466/silly",
		},
		{
			srcUrl:    "postgres://postgres:12345@localhost:5432/common?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			dbName:    "nh_common",
			resultUrl: "postgres://postgres:12345@localhost:5432/nh_common?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
		},
	}

	for _, c := range tCases {
		assert.Equal(t, ReplaceDbName(c.srcUrl, c.dbName), c.resultUrl)
	}
}
