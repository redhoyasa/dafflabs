package migration

import (
	"errors"
	"fmt"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

const (
	MigrationsPath = "migrations"
)

func ExecuteMigration(databaseUrl string) error {
	dir := fmt.Sprintf("file://%s", MigrationsPath)
	m, err := migrate.New(dir, databaseUrl)

	if err != nil {
		return fmt.Errorf("failed create migrate: %w", err)
	}

	m.Log = setupMigrationLogger()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate database error: %w", dbErr)
	}

	return nil
}

type logger struct {
	logger *log.Logger
}

func (l logger) Printf(arg string, vars ...interface{}) {
	l.logger.Printf(arg, vars...)
}

func (l logger) Verbose() bool {
	return true
}

func setupMigrationLogger() *logger {
	return &logger{
		logger: log.New(os.Stdout, "migrate", log.LstdFlags),
	}
}
