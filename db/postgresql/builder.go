package postgresql

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/igefined/go-kit/config"
)

const checkingSql = `select exists(select datname from pg_catalog.pg_database where datname = $1) as exist`

type QBuilder struct {
	pool *pgxpool.Pool

	logger *zap.Logger
}

func New(log *zap.Logger, cfg *config.DBCfg, lc fx.Lifecycle) *QBuilder {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if cfg.AutoCreateDatabase {
		CreateDatabase(ctx, log, cfg)
	}

	psqlCfg, err := pgx.ParseConfigWithOptions(cfg.URL, pgx.ParseConfigOptions{})
	if err != nil {
		log.Error("failed to parse postgres config", zap.Error(err))
	}

	conn, err := pgxpool.New(ctx, psqlCfg.ConnString())
	if err != nil {
		log.Error("failed to make postgres connection for auto create", zap.Error(err))
	}

	if lc != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				conn.Close()

				return nil
			},
			OnStart: func(ctx context.Context) error {
				if err = conn.Ping(ctx); err != nil {
					return fmt.Errorf("failed to ping database: %w", err)
				}

				return nil
			},
		})
	}

	return &QBuilder{conn, log}
}

func CreateDatabase(ctx context.Context, log *zap.Logger, cfg *config.DBCfg) {
	dbName := cfg.GetDatabaseName()
	conn, err := pgxpool.New(ctx, ReplaceDbName(cfg.URL, "postgres"))
	if err != nil {
		log.Error("failed to make postgres connection for auto create", zap.Error(err))
	}
	defer conn.Close()

	var exists bool

	row := conn.QueryRow(ctx, checkingSql, dbName)
	if err = row.Scan(&exists); err != nil {
		log.Error("autocreate db: failed to check the existence of the database", zap.Error(err))
	}

	if exists {
		log.Info("autocreate db: the database already exists", zap.String("db_name", dbName))
	} else {
		if _, err = conn.Exec(ctx, fmt.Sprintf(`create database "%s"`, dbName)); err != nil {
			log.Error("failed to create database", zap.Error(err))
		} else {
			log.Info("autocreate db: database created successfully")
		}
	}
}

func DropDatabase(ctx context.Context, log *zap.Logger, cfg *config.DBCfg) {
	dbName := cfg.GetDatabaseName()
	conn, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		log.Error("failed to make postgres connection for auto create", zap.Error(err))
	}
	defer conn.Close()

	var exists bool

	row := conn.QueryRow(ctx, checkingSql, dbName)
	if err = row.Scan(&exists); err != nil {
		log.Error("autocreate db: failed to check the existence of the database", zap.Error(err))
	}

	if exists {
		log.Info("autocreate db: the database already exists", zap.String("db_name", dbName))
	} else {
		if _, err := conn.Exec(ctx, fmt.Sprintf(`create database "%s"`, dbName)); err != nil {
			log.Error("failed to create database", zap.Error(err))
		} else {
			log.Info("autocreate db: database created successfully")
		}
	}
}

func ReplaceDbName(dbUrl, dbName string) string {
	parsed, err := url.Parse(dbUrl)
	if err != nil {
		return dbUrl
	}

	parsed.Path = "/" + dbName

	return parsed.String()
}

func (qb QBuilder) Querier() *pgxpool.Pool {
	return qb.pool
}

func (qb QBuilder) ConnString() string {
	return qb.pool.Config().ConnString()
}
