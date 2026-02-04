package app

import (
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) error {
	abs, err := filepath.Abs("./db/migrations")
	if err != nil {
		return err
	}
	migrationsPath := "file://" + filepath.ToSlash(abs)

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
