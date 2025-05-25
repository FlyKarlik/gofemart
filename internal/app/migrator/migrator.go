package migrator

import (
	"database/sql"
	"errors"

	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	ErrInvalidArgumnt = errors.New("usage up or down")
)

type AppMigrator struct {
	cfg     *config.Config
	logger  logger.Logger
	migrate *migrate.Migrate
}

func New(cfg *config.Config, logger logger.Logger) *AppMigrator {
	return &AppMigrator{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *AppMigrator) Migrate(migrateType string) error {
	db, err := sql.Open("postgres", a.cfg.Infra.Postgres.ConnStr)
	if err != nil {
		a.logger.Error("app[Migrator]", "AppMigrator.Migrate[sql.Open]", "Failed to create database connection", err)
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			a.logger.Error("app[Migrator]", "AppMigrator.Migrate[db.Close]", "Failed to close db connection", err)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		a.logger.Error("App", "AppMigrator.Migrate[postgres.WithInstance]", "Failed to create database driver", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+a.cfg.AppMigrator.MigrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		a.logger.Error("app[Migrator]", "AppMigrator.Migrate[migrate.NewDatabaseInstance]", "Failed to create migrator", err)
		return err
	}
	a.migrate = m

	switch migrateType {
	case "up":
		err = a.Up()
	case "down":
		err = a.Down()
	default:
		return ErrInvalidArgumnt
	}

	if err != nil {
		return err
	}

	return nil
}
