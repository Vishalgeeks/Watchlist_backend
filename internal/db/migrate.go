package db

import (
	"errors"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(connStr string) error {

	m, err := migrate.New(
		"file://./internal/db/migrations",
		connStr,
	)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to run")
			return nil
		}

		if strings.Contains(err.Error(), "Dirty database version") {
			log.Println("Dirty database detected, forcing version 8")
			m.Force(8)
			err = m.Up()
			if err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return err
			}
			return nil
		}

		return err
	}

	return nil
}
