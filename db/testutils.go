package db

import (
	"context"
	"database/sql"
	"os"

	"github.com/jackc/pgx/v5"
)

const DefaultTestDbURL = "postgres://postgres:postgres@localhost:5432/postgres"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func SetupTestConn() *pgx.Conn {
	dbURL := getEnv("PG_TEST_DB_URL", DefaultTestDbURL)
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = RunMigrations(db)
	if err != nil {
		panic(err)
	}
	err = RunFixtures(db)
	if err != nil {
		panic(err)
	}
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	return conn
}
