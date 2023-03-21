package cmd

import (
	"github.com/spf13/cobra"

	"github.com/nnkken/oracle-fetch/utils"
)

const FlagDatabaseURL = "database-url"

var RootCmd = &cobra.Command{
	Use:   "oracle-fetch",
	Short: "oracle-fetch is a tool to fetch data from the oracle",
	Long:  `oracle-fetch is a tool to fetch data from the oracle. It fetches data from the oracle and stores it into the database, and provide API to access the data.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.SetupLoggerFromCmdArgs(cmd)
	},
}

func init() {
	utils.ConfigCmd(RootCmd)
	RootCmd.PersistentFlags().String(FlagDatabaseURL, "postgres://postgres:postgres@localhost:5432/postgres", "database URL")

	RootCmd.AddCommand(FetchCmd)
	RootCmd.AddCommand(ApiCmd)
}
