package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/jackc/pgx/v5"

	"github.com/nnkken/oracle-fetch/datasource"
	"github.com/nnkken/oracle-fetch/datasource/types"
	"github.com/nnkken/oracle-fetch/db"
	"github.com/nnkken/oracle-fetch/runner"
)

const (
	FlagConfigFile    = "config"
	FlagFetchInterval = "fetch-interval"
)

func initDataSources(configFile string) ([]types.DataSource, error) {
	bz, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var configs []types.DataSourceConfig
	err = json.Unmarshal(bz, &configs)
	if err != nil {
		return nil, err
	}
	return datasource.InitDataSourcesFromConfig(configs)
}

var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch runs the service to fetch data from the oracle and insert them into the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbURL, err := cmd.Flags().GetString(FlagDatabaseURL)
		if err != nil {
			return err
		}
		fetchInterval, err := cmd.Flags().GetDuration(FlagFetchInterval)
		if err != nil {
			return err
		}
		chainLinkFile, err := cmd.Flags().GetString(FlagConfigFile)
		if err != nil {
			return err
		}

		conn, err := pgx.Connect(context.Background(), dbURL)
		if err != nil {
			return err
		}
		defer conn.Close(context.Background())

		dataSources, err := initDataSources(chainLinkFile)
		if err != nil {
			return err
		}

		database, err := sql.Open("pgx", dbURL)
		if err != nil {
			panic(err)
		}
		defer database.Close()
		err = db.RunMigrations(database)
		if err != nil {
			return err
		}
		database.Close()

		fetchLoop := runner.NewFetchLoop(fetchInterval)
		insertLoop := runner.NewInsertLoop()
		ch := make(chan db.DBEntry)

		go fetchLoop.Run(dataSources, ch)
		insertLoop.Run(conn, ch)
		return nil
	},
}

func init() {
	FetchCmd.Flags().String(FlagDatabaseURL, "postgres://postgres:postgres@localhost:5432/postgres", "Postgres database url")
	FetchCmd.Flags().Duration(FlagFetchInterval, 1*time.Minute, "fetch interval")
	FetchCmd.Flags().String(FlagConfigFile, "./config.json", "JSON file for data source definitions")
}
