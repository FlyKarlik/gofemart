package migrator

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
)

func (a *AppMigrator) Up() error {
	err := a.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			a.logger.Info("app[Migrator]", "AppMigrator.Up", "No migration changes found")
			return nil
		}

		a.logger.Error("app[Migrator]", "AppMigrator.Up[a.migrate.Up]", "Failed to migrate database", err)
		return err
	}

	a.logger.Info("app[Migrator]", "AppMigrator.Up", "Database migrated successfully")
	return nil
}
