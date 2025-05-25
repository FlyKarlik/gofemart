package main

import (
	"errors"
	"os"

	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/app/migrator"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	validator "github.com/FlyKarlik/gofemart/pkg/validation"

	_ "github.com/lib/pq"
)

var (
	ErrNotEnougthArgs = errors.New("not enouth arguments")
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic(ErrNotEnougthArgs)
	}

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := validator.Validate(cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.AppMigrator.LogLevel)
	if err != nil {
		panic(err)
	}

	migrator := migrator.New(cfg, logger)

	if err := migrator.Migrate(args[0]); err != nil {
		panic(err)
	}
}
