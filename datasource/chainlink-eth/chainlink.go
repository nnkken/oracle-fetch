package chainlink

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/nnkken/oracle-fetch/datasource"
	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth/contract"
	"github.com/nnkken/oracle-fetch/types"
)

var _ types.DataSource = (*ChainLinkETHSource)(nil)

type ChainLinkETHSource struct {
	instance *contract.Contract
	token    string
	unit     string
	decimals int
}

// NewChainLinkETHSource initialize the contract instance
func NewChainLinkETHSource(client *ethclient.Client, contractAddr common.Address, token, unit string, decimals int) (*ChainLinkETHSource, error) {
	instance, err := contract.NewContract(contractAddr, client)
	if err != nil {
		return nil, err
	}
	return &ChainLinkETHSource{
		instance: instance,
		token:    token,
		unit:     unit,
		decimals: decimals,
	}, nil
}

func (s *ChainLinkETHSource) Fetch() ([]types.DBEntry, error) {
	res, err := s.instance.LatestRoundData(nil)
	if err != nil {
		return nil, err
	}
	dbEntry := types.DBEntry{
		Token:          s.token,
		Unit:           s.unit,
		Price:          datasource.NormalizePrice(res.Answer, s.decimals),
		PriceTimestamp: time.Unix(res.UpdatedAt.Int64(), 0).UTC(),
		FetchTimestamp: time.Now().UTC(),
	}
	return []types.DBEntry{dbEntry}, nil
}
