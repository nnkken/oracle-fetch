package runner

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"go.uber.org/ratelimit"

	"github.com/nnkken/oracle-fetch/logger"
	"github.com/nnkken/oracle-fetch/types"
)

func RunFetchLoop(dataSources []types.DataSource, interval time.Duration, output chan<- types.DBEntry) {
	log := logger.GetLogger("fetch_loop")
	limiter := ratelimit.New(10)
	timech := time.Tick(interval)
	for ; true; <-timech {
		for _, datasource := range dataSources {
			go func(datasource types.DataSource) {
				// TODO: move limiter into ethclient
				limiter.Take()
				dbentries, err := datasource.Fetch()
				if err != nil {
					log.Errorw("failed to fetch data", "error", err)
					return
				}
				for _, dbentry := range dbentries {
					output <- dbentry
				}
			}(datasource)
		}
	}
}

func RunInsertLoop(conn *pgx.Conn, ch <-chan types.DBEntry) {
	log := logger.GetLogger("insert_loop")
	for entry := range ch {
		_, err := conn.Exec(context.Background(), `
			INSERT INTO prices (token, unit, price, price_timestamp, fetch_timestamp)
			VALUES ($1, $2, $3, $4, $5)
		`, entry.Token, entry.Unit, entry.Price, entry.PriceTimestamp, entry.FetchTimestamp)
		if err != nil {
			log.Errorw("failed to insert entry into database", "error", err, "entry", entry)
		}
		log.Infow("inserted entry into database", "entry", entry)
	}
}

func Run(dataSources []types.DataSource, interval time.Duration, conn *pgx.Conn) {
	ch := make(chan types.DBEntry)
	go RunFetchLoop(dataSources, interval, ch)
	RunInsertLoop(conn, ch)
}
