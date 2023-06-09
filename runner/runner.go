package runner

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/nnkken/oracle-fetch/datasource/types"
	"github.com/nnkken/oracle-fetch/db"
	"github.com/nnkken/oracle-fetch/utils"
)

type FetchLoop struct {
	interval time.Duration
}

func NewFetchLoop(interval time.Duration) *FetchLoop {
	return &FetchLoop{
		interval: interval,
	}
}

func fetch(dataSource types.DataSource, output chan<- db.DBEntry, log *zap.SugaredLogger) {
	dbEntries, err := dataSource.Fetch()
	if err != nil {
		log.Errorw("failed to fetch data", "error", err)
		return
	}
	for _, dbEntry := range dbEntries {
		output <- dbEntry
	}
}

func (loop *FetchLoop) Run(dataSources []types.DataSource, output chan<- db.DBEntry) {
	log := utils.GetLogger("fetch_loop")
	timeCh := time.Tick(loop.interval)
	// not using range since we want to execute loop immediately before the first tick
	for ; true; <-timeCh {
		for _, dataSource := range dataSources {
			go fetch(dataSource, output, log)
		}
	}
}

type InsertLoop struct{}

func NewInsertLoop() *InsertLoop {
	return &InsertLoop{}
}

func Insert(entry db.DBEntry, conn *pgx.Conn, log *zap.SugaredLogger) {
	_, err := conn.Exec(context.Background(), `
			INSERT INTO prices (token, unit, price, price_timestamp, fetch_timestamp, source)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, entry.Token, entry.Unit, entry.Price, entry.PriceTimestamp, entry.FetchTimestamp, entry.Source)
	if err != nil {
		log.Errorw("failed to insert entry into database", "error", err, "entry", entry)
		return
	}
	log.Infow("inserted entry into database", "entry", entry)
}

func (loop *InsertLoop) Run(conn *pgx.Conn, ch <-chan db.DBEntry) {
	log := utils.GetLogger("insert_loop")
	for entry := range ch {
		Insert(entry, conn, log)
	}
}
