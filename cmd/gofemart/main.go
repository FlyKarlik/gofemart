package main

import (
	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/app/gofemart"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	validator "github.com/FlyKarlik/gofemart/pkg/validation"
)

// @title GoFemart API
// @version 1.0
// @description API documentation for the GoFemart backend service.

// @contact.name API Support
// @contact.url https://github.com/FlyKarlik/gofemart
// @contact.email nikitasavin191@gmail.com

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := validator.Validate(cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.AppGofemart.LogLevel)
	if err != nil {
		panic(err)
	}

	gofemart := gofemart.New(cfg, logger)
	if err := gofemart.Start(); err != nil {
		panic(err)
	}
}
