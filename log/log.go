package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level) (*zap.Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

func WithAppInfo(logger *zap.Logger, version, commit, buildDate string) {
	logger.With(
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("buildDate", buildDate),
	)
}
