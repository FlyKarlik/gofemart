package server

import (
	"context"
	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/handler"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/router"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

type HTTPServer struct {
	cfg        *config.Config
	router     *router.HTTPRouter
	handler    *handler.Handler
	httpserver *http.Server
}

func New(
	cfg *config.Config,
	logger logger.Logger,
	router *router.HTTPRouter,
	handler *handler.Handler) *HTTPServer {

	httpServer := &HTTPServer{
		cfg:     cfg,
		router:  router,
		handler: handler,
	}

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.AppGofemart.AppHost, cfg.AppGofemart.AppPort),
		Handler:        router.InitRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	httpServer.httpserver = srv

	return httpServer
}

func (h *HTTPServer) ListenAndServe() error {
	return h.httpserver.ListenAndServe()
}

func (h *HTTPServer) Shuttdown(ctx context.Context) error {
	return h.httpserver.Shutdown(ctx)
}
