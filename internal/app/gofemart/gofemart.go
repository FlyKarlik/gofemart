package gofemart

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/handler"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/middleware"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/router"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/server"
	"github.com/FlyKarlik/gofemart/internal/repository"
	"github.com/FlyKarlik/gofemart/internal/usecase"
	"github.com/FlyKarlik/gofemart/pkg/database"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"github.com/FlyKarlik/gofemart/pkg/trace"
)

type AppGofemart struct {
	cfg    *config.Config
	logger logger.Logger
}

func New(cfg *config.Config, logger logger.Logger) *AppGofemart {
	return &AppGofemart{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *AppGofemart) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.logger.Info("app[Gofemart]", "AppGofemart.Start", "Gofemart application started...")

	shuttdownTrace, err := trace.New(context.Background(), &a.cfg.Infra.Jaeger)
	if err != nil {
		a.logger.Error("app[Gofemart]", "AppGofemart.Start[trace.New]", "Failed to init jaeger tracer", err)
		return err
	}
	defer func() {
		if err := shuttdownTrace(context.Background()); err != nil {
			a.logger.Error("app[Gofemart]", "AppGofemart.Start[shuttdownTrace]", "Failed to shuttdown tracer", err)
		}
	}()

	postgresConn, err := database.NewPostgresDB(&a.cfg.Infra.Postgres)
	if err != nil {
		a.logger.Error("app[Gofemart]", "AppGofemart.Start[database.NewPostgresDB]", "Failed to init postgresql", err)
		return err
	}
	defer postgresConn.Close()

	redisClient := database.NewRedisClient(&a.cfg.Infra.Redis)

	repo := repository.New(a.logger, postgresConn, redisClient)
	usecase := usecase.New(a.cfg, a.logger, repo)

	httpHandler := handler.New(a.logger, usecase)
	httpMiddleware := middleware.New(a.cfg, a.logger, usecase)
	httpRouter := router.New(httpMiddleware, httpHandler)
	httpServer := server.New(a.cfg, a.logger, httpRouter, httpHandler)

	go func() {
		a.logger.Infof(
			"app[Gofemart]",
			"AppMigrator.Start[httpServer.ListenAndServe]",
			"HTTP server starting...",
			"host: %s, port: %s", a.cfg.AppGofemart.AppHost, a.cfg.AppGofemart.AppPort)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("app[Gofemart]", "AppMigrator.Start[httpServer.ListenAndServer]", "Failed while started http server", err)
			os.Exit(1)
		}
	}()

	a.signalHandler(ctx)
	if err := httpServer.Shuttdown(ctx); err != nil {
		a.logger.Error("app[Gofemart]", "AppMigrator.Start[httpServer.Shuttdown]", "Failed while shuttdown http server", err)
		return err
	}

	return nil
}

func (a *AppGofemart) signalHandler(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		signal.Stop(signalChan)
		close(signalChan)
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("app[Gofermat]", "AppMigrator.signalHandler", "Context cancelled")
		return
	case <-signalChan:
		a.logger.Info("app[Gofemart]", "AppMigrator.signalHandler", "Signal received")
		return
	}
}
