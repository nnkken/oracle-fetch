package db

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(dbURL string) (err error) {
	goose.SetBaseFS(embedMigrations)
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}
	return nil
}
