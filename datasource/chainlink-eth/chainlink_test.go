package chainlink

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/nnkken/oracle-fetch/db"
	"github.com/nnkken/oracle-fetch/utils"
)

var _ ChainLinkContract = (*MockChainLinkContract)(nil)

func TestChainLinkETHSource(t *testing.T) {
	utils.MockTimeNow(t, time.Unix(1300000000, 0))

	token := "TOKEN"
	unit := "UNIT"
	decimals := uint8(8)
	mockInstance := NewMockChainLinkContract(t)

	source := NewChainLinkETHSource(mockInstance, token, unit, decimals)

	mockCall := mockInstance.EXPECT().LatestRoundData((*bind.CallOpts)(nil)).Return(RoundData{
		Answer:    big.NewInt(1e10),
		StartedAt: big.NewInt(1000000000),
		UpdatedAt: big.NewInt(1234567890),
	}, nil)

	dbEntries, err := source.Fetch()
	require.NoError(t, err)
	require.Len(t, dbEntries, 1)
	require.Equal(t, db.DBEntry{
		Token:          token,
		Unit:           unit,
		Price:          big.NewInt(1e10),
		PriceTimestamp: time.Unix(1234567890, 0).UTC(),
		FetchTimestamp: utils.TimeNow().UTC(),
		Source:         "chainlink-eth",
	}, dbEntries[0])
	mockCall.Unset()

	expectedErr := errors.New("err")
	mockCall = mockInstance.EXPECT().LatestRoundData((*bind.CallOpts)(nil)).Return(RoundData{}, expectedErr)
	_, err = source.Fetch()
	require.ErrorIs(t, err, expectedErr)
	mockCall.Unset()

	decimals = uint8(10)
	source = NewChainLinkETHSource(mockInstance, token, unit, decimals)
	mockCall = mockInstance.EXPECT().LatestRoundData((*bind.CallOpts)(nil)).Return(RoundData{
		Answer:    big.NewInt(1e10),
		StartedAt: big.NewInt(1000000000),
		UpdatedAt: big.NewInt(1234567890),
	}, nil)
	dbEntries, err = source.Fetch()
	require.NoError(t, err)
	require.Len(t, dbEntries, 1)
	require.Equal(t, db.DBEntry{
		Token:          token,
		Unit:           unit,
		Price:          big.NewInt(1e8),
		PriceTimestamp: time.Unix(1234567890, 0).UTC(),
		FetchTimestamp: utils.TimeNow().UTC(),
		Source:         "chainlink-eth",
	}, dbEntries[0])
	mockCall.Unset()

	decimals = uint8(6)
	source = NewChainLinkETHSource(mockInstance, token, unit, decimals)
	mockCall = mockInstance.EXPECT().LatestRoundData((*bind.CallOpts)(nil)).Return(RoundData{
		Answer:    big.NewInt(1e10),
		StartedAt: big.NewInt(1000000000),
		UpdatedAt: big.NewInt(1234567890),
	}, nil)
	dbEntries, err = source.Fetch()
	require.NoError(t, err)
	require.Len(t, dbEntries, 1)
	require.Equal(t, db.DBEntry{
		Token:          token,
		Unit:           unit,
		Price:          big.NewInt(1e12),
		PriceTimestamp: time.Unix(1234567890, 0).UTC(),
		FetchTimestamp: utils.TimeNow().UTC(),
		Source:         "chainlink-eth",
	}, dbEntries[0])
	mockCall.Unset()
}
