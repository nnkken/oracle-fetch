package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nnkken/oracle-fetch/types"

	"go.uber.org/ratelimit"
)

func RunFetchLoop(dataSources []types.DataSource, interval time.Duration, output chan<- types.DBEntry) {
	limiter := ratelimit.New(10)
	timech := time.Tick(interval)
	for ; true; <-timech {
		for _, datasource := range dataSources {
			go func(datasource types.DataSource) {
				// TODO: move limiter into ethclient
				limiter.Take()
				dbentries, err := datasource.Fetch()
				if err != nil {
					// TODO: log
					return
				}
				for _, dbentry := range dbentries {
					output <- dbentry
				}
			}(datasource)
		}
	}
}

func RunInsertLoop(pool *pgxpool.Pool, ch <-chan types.DBEntry) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	defer conn.Release()

	// TODO: use some migration tool to manage this
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS prices (
			id SERIAL PRIMARY KEY,
			token TEXT,
			unit TEXT,
			price NUMERIC,
			price_timestamp TIMESTAMP WITH TIME ZONE,
			fetch_timestamp TIMESTAMP WITH TIME ZONE
		)
	`)
	if err != nil {
		panic(err)
	}
	for entry := range ch {
		now := time.Now().UTC().Format(time.RFC3339)
		_, err := conn.Exec(context.Background(), `
			INSERT INTO prices (token, unit, price, price_timestamp, fetch_timestamp)
			VALUES ($1, $2, $3, $4, $5)
		`, entry.Token, entry.Unit, entry.Price, entry.PriceTimestamp, entry.FetchTimestamp)
		if err != nil {
			// TODO: non-critical, log error and ignore instead of panic
			panic(err)
		}
		fmt.Printf("%s - %##v\n", now, entry)
	}
}

func Run(dataSources []types.DataSource, interval time.Duration, pool *pgxpool.Pool) {
	ch := make(chan types.DBEntry)
	go RunFetchLoop(dataSources, interval, ch)
	RunInsertLoop(pool, ch)
}
