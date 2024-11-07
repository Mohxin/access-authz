package authz

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/volvo-cars/connect-access-control/internal/api"
	v1 "github.com/volvo-cars/connect-access-control/internal/api/v1"
	"github.com/volvo-cars/connect-access-control/internal/config"
	"github.com/volvo-cars/connect-access-control/internal/pkg/authz"
	cachemanager "github.com/volvo-cars/connect-access-control/internal/pkg/gateway/cache-manager"
	"github.com/volvo-cars/connect-access-control/internal/pkg/gateway/plums"
	"github.com/volvo-cars/connect-access-control/internal/pkg/store"
	httpserver "github.com/volvo-cars/go-ecp-httpserver"
	"github.com/volvo-cars/go-middlewares"
	"github.com/volvo-cars/go-observer"
	"github.com/volvo-cars/go-tracer/otel"
)

// @title			Access Control API
// @version		1.0
// @description	Access Control API for the Volvo Cars Connect platform.
// @BasePath		/v1
func Run(cfg *config.Config) {
	// Set up OpenTelemetry :: create tracer
	tp, err := NewTracerProvider(context.Background(), cfg)
	if err != nil {
		slog.Error("failed to start open telemetry tracer", slog.Any("error", err))
		return
	}

	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown open telemetry tracer", slog.Any("error", err))
		}
	}()

	store := store.NewAccessControlStore(cfg.IAM.RootDir)
	if err = store.Process(); err != nil {
		slog.Error("failed to load access-control in-memory data", slog.Any("error", err))
		return
	}

	plumsCfg, err := plums.LoadConfig()
	if err != nil {
		slog.Error("failed to load plums config", slog.Any("error", err))
		return
	}

	cacheManagerCfg, err := cachemanager.LoadConfig()
	if err != nil {
		slog.Error("failed to load cache-manager config", slog.Any("error", err))
		return
	}

	outgoingCollector := observer.NewOutgoingCollector(cfg.App.Name)
	plumsClient := plums.New(plumsCfg, outgoingCollector)
	cacheManagerClient := cachemanager.New(cacheManagerCfg, outgoingCollector)
	authClient := authz.NewService(cacheManagerClient, plumsClient, store)

	// main router
	r := chi.NewRouter()

	// swagger
	if cfg.App.Environment != config.EnvironmentProduction {
		r.Mount("/swagger", httpSwagger.WrapHandler)
	}

	// controllers
	controllers := []api.Controller{
		v1.NewController(store, authClient),
	}

	r.Mount("/v1", api.RegisterRoutes(NewAPIRouter(cfg), controllers...))

	// http server
	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))
	slog.Info("http server listening", slog.Any("port", cfg.HTTP.Port), slog.Any("environment", cfg.App.Environment.String()), slog.Any("App Name", cfg.App.Name))

	// Graceful shutdown
	// Create a channel that listens on incoming interrupt signals
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case err := <-httpServer.Notify():
		slog.Error("http server failed to start serve", slog.Any("error", slog.Any("error", err)))
	case sig := <-exit:
		slog.Info("http server received termination signal", slog.Any("signal", sig))
	}

	if err = httpServer.Shutdown(); err != nil {
		slog.Error("http server failed to shutdown", slog.Any("error", slog.Any("error", err)))
		return
	}

	slog.Info("http server shutdown successfully")
}

func NewAPIRouter(cfg *config.Config) *chi.Mux {
	observerMiddleware := observer.New(observer.NewMetricCollector(cfg.App.Name))

	router := chi.NewRouter()

	// These middlewares (RealIP, RequestId, CorrelationId) should be inserted fairly early in the middleware stack to
	// ensure that subsequent layers (e.g., request loggers)
	router.Use(middleware.RealIP)
	router.Use(middlewares.RequestId)
	router.Use(middlewares.CorrelationId)
	router.Use(middleware.NoCache)
	router.Use(otel.Handler)
	router.Use(observerMiddleware.Handler)
	router.Use(middlewares.RequestLogger())
	router.Use(middleware.Recoverer)
	return router
}

func NewTracerProvider(ctx context.Context, cfg *config.Config) (*otel.TracerProvider, error) {
	return otel.NewTracerProvider(
		ctx,
		otel.WithServiceName(cfg.App.Name),
		otel.WithVersion(cfg.App.Version),
		otel.WithEnvironment(cfg.App.Environment.String()),
		otel.WithHTTPEndpointURL(cfg.Tracer.EndpointURL),
	)
}
