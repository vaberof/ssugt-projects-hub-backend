package main

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"ssugt-projects-hub/config"
)

func runMigrations(db *sqlx.DB) (err error) {
	err = runPostgresMigrations(db)
	if err != nil {
		return
	}

	return nil
}

func runPostgresMigrations(db *sqlx.DB) (err error) {
	return goose.Up(db.DB, "database/postgres/migrations")
}

func runPostgresMigrations2() (err error) {
	postgres.DefaultMigrationsTable = "ssugt_projects_hub_migrations"
	m, err := migrate.New(
		"file://database/postgres/migrations",
		config.PostgresConnection())
	if err != nil {
		return
	}

	err = m.Up()
	if err != nil {
		if err.Error() == migrate.ErrNoChange.Error() {
			return nil
		}
		return
	}

	return nil
}
