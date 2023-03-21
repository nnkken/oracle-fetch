package chainlink

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth/contract"
	"github.com/nnkken/oracle-fetch/datasource/types"
)

type ChainLinkEthContractEntry struct {
	Token    string `json:"token"`
	Unit     string `json:"unit"`
	Decimals uint8  `json:"decimals"`
	Address  string `json:"address"`
}

type ChainLinkEthConfig struct {
	EthEndpoint string                      `json:"eth_endpoint"`
	Contracts   []ChainLinkEthContractEntry `json:"contracts"`
}

func NewDataSourceFromConfig(rawConfig json.RawMessage) ([]types.DataSource, error) {
	var config ChainLinkEthConfig
	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(config.EthEndpoint)
	if err != nil {
		return nil, err
	}
	// TODO: handle client close

	dataSources := make([]types.DataSource, len(config.Contracts))
	for i, entry := range config.Contracts {
		contractInstance, err := contract.NewContract(common.HexToAddress(entry.Address), client)
		if err != nil {
			return nil, err
		}
		dataSources[i] = NewChainLinkETHSource(contractInstance, entry.Token, entry.Unit, entry.Decimals)
	}
	return dataSources, nil
}
