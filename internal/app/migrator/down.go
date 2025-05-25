package migrator

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
)

func (a *AppMigrator) Down() error {
	err := a.migrate.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			a.logger.Info("app[Migrator]", "AppMigrator.Down", "No down migrations to apply")
			return nil
		}

		a.logger.Error("app[Migrator]", "AppMigrator.Down[a.migrate.Down]", "Failed to migrate database", err)
		return err
	}

	a.logger.Info("app[Migrator]", "AppMigrator.Down", "Database migrated successfully")
	return nil
}
