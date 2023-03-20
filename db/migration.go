package db

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}
	return nil
}
