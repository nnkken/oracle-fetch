package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nnkken/oracle-fetch/api"
)

const FlagHost = "host"

var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbURL, err := cmd.Flags().GetString(FlagDatabaseURL)
		if err != nil {
			return err
		}
		host, err := cmd.Flags().GetString(FlagHost)
		if err != nil {
			return err
		}

		connPool, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			return err
		}
		defer connPool.Close()

		return api.NewRouter(connPool).Run(host)
	},
}

func init() {
	ApiCmd.Flags().String(FlagHost, "localhost:8080", "host address")
}
