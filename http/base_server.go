package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/igdotog/core/config"
	"github.com/igdotog/core/logger"

	"go.uber.org/fx"
)

type BaseServer struct {
	Mux *chi.Mux

	cfg *config.Application
	log *logger.Logger
}

func NewBaseServer(
	cfg *config.Application,
	log *logger.Logger,
	lc fx.Lifecycle,
) *BaseServer {
	router := chi.NewRouter()

	httpServer := http.Server{
		Addr:        ":" + strconv.Itoa(int(cfg.Http.Port)),
		IdleTimeout: 30 * time.Second,
		Handler:     router,
	}

	baseServer := &BaseServer{
		Mux: router,
		cfg: cfg,
		log: log,
	}

	serve := func() {
		baseServer.Mux.Group(
			func(statusGrp chi.Router) {
				statusGrp.Handle("/status", HealthHandler())
			},
		)

		httpServer.Handler = baseServer.Mux

		log.Infof("The base http-server listen and serve on %d port", cfg.Http.Port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}

	if lc != nil {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go serve()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := httpServer.Shutdown(ctx); err != http.ErrServerClosed {
					return err
				}

				return nil
			},
		})
	} else {
		go serve()
	}

	return baseServer
}

// HealthHandler godoc
// @summary Health-Check
// @success 200 {string} string "ok"
// @tags Service status
// @router /status [get]
func HealthHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	}
}
