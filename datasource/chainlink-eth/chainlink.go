package chainlink

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth/contract"
	"github.com/nnkken/oracle-fetch/datasource/types"
	"github.com/nnkken/oracle-fetch/db"
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
	Token    string
	Unit     string
	Decimals uint8
}

// NewChainLinkETHSource initialize the contract instance
func NewChainLinkETHSource(instance ChainLinkContract, token, unit string, decimals uint8) *ChainLinkETHSource {
	return &ChainLinkETHSource{
		instance: instance,
		Token:    token,
		Unit:     unit,
		Decimals: decimals,
	}
}

func (s *ChainLinkETHSource) Fetch() ([]db.DBEntry, error) {
	res, err := s.instance.LatestRoundData(nil)
	if err != nil {
		return nil, err
	}
	dbEntry := db.DBEntry{
		Token:          s.Token,
		Unit:           s.Unit,
		Price:          utils.NormalizePrice(res.Answer, s.Decimals),
		PriceTimestamp: time.Unix(res.UpdatedAt.Int64(), 0).UTC(),
		FetchTimestamp: utils.TimeNow().UTC(),
	}
	return []db.DBEntry{dbEntry}, nil
}
