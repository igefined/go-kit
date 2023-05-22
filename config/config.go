package config

import (
	"context"
	"log"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

type BaseConfig struct {
	Name  string `config:"APP_NAME"`
	Port  uint   `config:"PORT"`
	DBUrl string `config:"DB_URL,required"`
}

func NewBaseConfig(ctx context.Context) *BaseConfig {
	cfgCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	cfg := &BaseConfig{}
	loader := confita.NewLoader(env.NewBackend())
	err := loader.Load(cfgCtx, cfg)
	if err != nil {
		log.Fatal("config has not loaded")
	}

	return cfg
}
