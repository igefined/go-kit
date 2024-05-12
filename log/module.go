package log

import (
	"context"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Options(
	fx.WithLogger(func(logger *zap.SugaredLogger) fxevent.Logger { return &fxevent.ZapLogger{Logger: logger.Desugar()} }),
	fx.Provide(
		func(ls fx.Lifecycle) (*Logger, error) {
			logLevel := zapcore.InfoLevel
			logLevelStr := os.Getenv("LOG_LEVEL")
			if logLevelStr != "" {
				if err := logLevel.UnmarshalText([]byte(logLevelStr)); err != nil {
					return nil, err
				}
			}

			l, err := NewLogger(logLevel)
			if err != nil {
				return nil, err
			}

			ls.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					_ = l.Sync()
					return nil
				},
				OnStop: func(_ context.Context) error {
					return nil
				},
			})

			return l, err
		},

		func(z *zap.Logger) *zap.SugaredLogger {
			return z.Sugar()
		},
	),
)
