package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/jackc/pgx/v5"

	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth"
	"github.com/nnkken/oracle-fetch/db"
	"github.com/nnkken/oracle-fetch/runner"
	"github.com/nnkken/oracle-fetch/types"
)

const (
	FlagChainLinkFile = "chain-link-file"
	FlagFetchInterval = "fetch-interval"
	FlagEthEndpoint   = "eth-endpoint"
)

type ChainLinkJsonEntry struct {
	Token    string `json:"token"`
	Unit     string `json:"unit"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
}

func initChainLinkETH(chainLinkFile string, client *ethclient.Client) ([]types.DataSource, error) {
	bz, err := os.ReadFile(chainLinkFile)
	if err != nil {
		panic(err)
	}
	var entries []ChainLinkJsonEntry
	err = json.Unmarshal(bz, &entries)
	if err != nil {
		panic(err)
	}

	dataSources := make([]types.DataSource, len(entries))
	for i, entry := range entries {
		dataSource, err := chainlink.NewChainLinkETHSource(client, common.HexToAddress(entry.Address), entry.Token, entry.Unit, entry.Decimals)
		if err != nil {
			panic(err)
		}
		dataSources[i] = dataSource
	}
	return dataSources, nil
}

var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch runs the service to fetch data from the oracle and insert them into the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbURL, err := cmd.Flags().GetString(FlagDatabaseURL)
		if err != nil {
			return err
		}
		ethEndpoint, err := cmd.Flags().GetString(FlagEthEndpoint)
		if err != nil {
			return err
		}
		fetchInterval, err := cmd.Flags().GetDuration(FlagFetchInterval)
		if err != nil {
			return err
		}
		chainLinkFile, err := cmd.Flags().GetString(FlagChainLinkFile)
		if err != nil {
			return err
		}

		conn, err := pgx.Connect(context.Background(), dbURL)
		if err != nil {
			return err
		}
		defer conn.Close(context.Background())

		// TODO: integrate rate limiter
		client, err := ethclient.Dial(ethEndpoint)
		if err != nil {
			return err
		}
		defer client.Close()

		dataSources, err := initChainLinkETH(chainLinkFile, client)
		if err != nil {
			return err
		}

		err = db.RunMigrations(dbURL)
		if err != nil {
			return err
		}
		runner.Run(dataSources, fetchInterval, conn)
		return nil
	},
}

func init() {
	FetchCmd.Flags().String(FlagDatabaseURL, "postgres://postgres:postgres@localhost:5432/postgres", "Postgres database url")
	FetchCmd.Flags().String(FlagEthEndpoint, "http://localhost:5051", "Ethereum endpoint")
	FetchCmd.Flags().Duration(FlagFetchInterval, 1*time.Minute, "fetch interval")
	FetchCmd.Flags().String(FlagChainLinkFile, "./chain-link.json", "ChainLink JSON file for address and token info")
}
