package config

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

const (
	defaultHttpPort = 8080
	defaultAppName  = "app"
	defaultEnv      = "DEV"
)

type (
	Application struct {
		Name  string `config:"APP_NAME"`
		Env   string `config:"APP_ENV"`
		AppId string

		*DB
		*Http
	}

	Http struct {
		Port uint `config:"HTTP_PORT"`
	}

	DB struct {
		URL                string `config:"DB_URL"`
		AutoCreateDatabase bool   `config:"DB_AUTO_CREATE_DATABASE"`
	}
)

func New(ctx context.Context, other ...interface{}) (*Application, error) {
	cfg := Application{
		Name: defaultAppName,
		Env:  defaultEnv,
	}

	for _, otherCfg := range other {
		switch v := otherCfg.(type) {
		case *Http:
			if v.Port == 0 {
				v.Port = defaultHttpPort
			}
			cfg.Http = v
		case *DB:
			cfg.DB = v
		default:
			return nil, fmt.Errorf("failed to load config, %T is not undefined type", otherCfg)
		}
	}

	err := confita.NewLoader(env.NewBackend()).Load(ctx, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %s", err)
	}

	cfg.AppId = CreateAppID(cfg.Name)

	return &cfg, nil
}

func CreateAppID(prefix string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(prefix), hex.EncodeToString(uuid.NodeID()))
}
