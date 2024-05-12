package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/igefined/go-kit/refid"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(level zapcore.Level) (*Logger, error) {
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

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger}, nil
}

func (l *Logger) M(ctx context.Context) *Logger {
	return l.With(Markers(ctx)...)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

func Markers(ctx context.Context) []zap.Field {
	tags := refid.GetFields(ctx)
	fields := make([]zap.Field, 0, len(tags))
	for k, v := range tags {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func WithAppInfo(logger *zap.Logger, version, commit, buildDate string) {
	logger.With(
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("buildDate", buildDate),
	)
}
