package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Папка с миграциями
		"postgres",          // Название базы данных
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("[INFO] Migrations applied successfully")
	return nil
}

func RollbackMigrations(db *sql.DB, steps int) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Printf("[INFO] Rolled back %d migrations", steps)
	return nil
}

func GetMigrationVersion(db *sql.DB) (uint, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return 0, err
	}

	version, dirty, err := m.Version()
	if err != nil {
		return 0, err
	}

	if dirty {
		log.Println("[WARN] The migration is in a dirty state")
	}

	return version, nil
}
