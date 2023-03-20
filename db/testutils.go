package db

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const DefaultTestDbURL = "postgres://postgres:postgres@localhost:5432/postgres"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func SetupTestConn(t *testing.T) *pgx.Conn {
	dbURL := getEnv("PG_TEST_DB_URL", DefaultTestDbURL)
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	rows, err := conn.Query(context.Background(), "SELECT 1 FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		panic(errors.New("testing database is not empty, possibly testing on production database"))
	}
	err = RunMigrations(dbURL)
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		conn.Close(context.Background())
		conn, err := pgx.Connect(context.Background(), dbURL)
		if err != nil {
			panic(err)
		}
		defer conn.Close(context.Background())
		var tables []string
		rows, err := conn.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			var tableName string
			err := rows.Scan(&tableName)
			if err != nil {
				panic(err)
			}
			tables = append(tables, tableName)
		}
		rows.Close() // for freeing the connection; it's OK to duplicate close as stated in pgx docs
		for _, tableName := range tables {
			_, err := conn.Exec(context.Background(), "DROP TABLE "+tableName)
			if err != nil {
				panic(err)
			}
		}
		conn.Close(context.Background())
	})
	return conn
}
