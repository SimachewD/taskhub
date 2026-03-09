package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func RunMigrations(dsn string) error {
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return err
	}
	defer m.Close() // ensures connections are released

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}