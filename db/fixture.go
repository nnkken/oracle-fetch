package db

import (
	"database/sql"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/go-testfixtures/testfixtures/v3"
)

//go:embed fixtures/*.yml
var embedFixtures embed.FS

func RunFixtures(db *sql.DB) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.FS(embedFixtures),
		testfixtures.Directory("fixtures"),
	)
	if err != nil {
		return err
	}
	err = fixtures.EnsureTestDatabase()
	if err != nil {
		return err
	}
	err = fixtures.Load()
	if err != nil {
		return err
	}
	return nil
}
