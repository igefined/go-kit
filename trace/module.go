package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.uber.org/fx"

	"github.com/igefined/go-kit/config"
)

func Module(serviceName string) fx.Option {
	return fx.Options(
		fx.Provide(),
		fx.Invoke(func(
			ls fx.Lifecycle,
			appCfg *config.MainCfg,
			traceCfg *config.TraceCfg,
			client otlptrace.Client,
		) {
			ls.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					exporter, err := NewTraceSpanExporter(ctx, client)
					if err != nil {
						return err
					}

					NewGlobalTracerProvider(serviceName, exporter)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)
}
