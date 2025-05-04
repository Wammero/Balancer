package migration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func callMigrations(dsn string) (*migrate.Migrate, error) {
	migrationsPath := "file://migrations"

	m, err := migrate.New(migrationsPath, dsn)

	if err != nil {
		return nil, err
	}

	return m, err
}

func ApplyMigrations(connStr string) {
	m, err := callMigrations(connStr)
	if err != nil {
		logrus.Fatalf("Error initializing migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			logrus.Info("Migrations are already applied")
		} else {
			logrus.Fatalf("Error applying migrations: %v", err)
		}
	} else {
		logrus.Info("Migrations successfully applied")
	}
}
