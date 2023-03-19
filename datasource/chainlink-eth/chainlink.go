package chainlink

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/nnkken/oracle-fetch/datasource"
	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth/contract"
	"github.com/nnkken/oracle-fetch/types"
	"github.com/nnkken/oracle-fetch/utils"
)

var (
	_ types.DataSource  = (*ChainLinkETHSource)(nil)
	_ ChainLinkContract = (*contract.Contract)(nil)
)

type RoundData = struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

//go:generate --name ChainLinkContractMock
type ChainLinkContract interface {
	LatestRoundData(opts *bind.CallOpts) (RoundData, error)
}

type ChainLinkETHSource struct {
	instance ChainLinkContract
	token    string
	unit     string
	decimals uint8
}

// NewChainLinkETHSource initialize the contract instance
func NewChainLinkETHSource(instance ChainLinkContract, token, unit string, decimals uint8) *ChainLinkETHSource {
	return &ChainLinkETHSource{
		instance: instance,
		token:    token,
		unit:     unit,
		decimals: decimals,
	}
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
		FetchTimestamp: utils.TimeNow().UTC(),
	}
	return []types.DBEntry{dbEntry}, nil
}
