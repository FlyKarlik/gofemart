package handler

import (
	"github.com/FlyKarlik/gofemart/internal/usecase"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"time"
)

type Handler struct {
	logger  logger.Logger
	usecase *usecase.Usecase
	startup time.Time
}

func New(logger logger.Logger, usecase *usecase.Usecase) *Handler {
	return &Handler{
		startup: time.Now(),
		logger:  logger,
		usecase: usecase,
	}
}
