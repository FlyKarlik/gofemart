package middleware

import (
	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/usecase"
	"github.com/FlyKarlik/gofemart/pkg/logger"
)

type Middleware struct {
	cfg     *config.Config
	logger  logger.Logger
	usecase *usecase.Usecase
}

func New(cfg *config.Config, logger logger.Logger, usecase *usecase.Usecase) *Middleware {
	return &Middleware{
		cfg:     cfg,
		logger:  logger,
		usecase: usecase}
}
