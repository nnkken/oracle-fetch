package datasource

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth"
	"github.com/nnkken/oracle-fetch/types"
)

func TestInitDataSourcesFromConfig(t *testing.T) {
	config := []types.DataSourceConfig{
		{
			Type: "chainlink-eth",
			Config: []byte(`{
				"eth_endpoint": "http://somewhere",
				"contracts": [
					{
						"token": "TOKEN-0",
						"unit": "UNIT-0",
						"decimals": 8,
						"address": "0x0000000000000000000000000000000000000000"
					},
					{
						"token": "TOKEN-1",
						"unit": "UNIT-1",
						"decimals": 10,
						"address": "0x1111111111111111111111111111111111111111"
					}
				]
			}`),
		},
		{
			Type:         "chainlink-eth",
			RateLimitRps: 10,
			Config: []byte(`{
				"eth_endpoint": "http://somewhere2",
				"contracts": [
					{
						"token": "TOKEN-2",
						"unit": "UNIT-2",
						"decimals": 8,
						"address": "0x2222222222222222222222222222222222222222"
					},
					{
						"token": "TOKEN-3",
						"unit": "UNIT-3",
						"decimals": 10,
						"address": "0x3333333333333333333333333333333333333333"
					}
				]
			}`),
		},
	}
	dataSources, err := InitDataSourcesFromConfig(config)
	require.NoError(t, err)
	require.Len(t, dataSources, 4)

	require.IsType(t, &chainlink.ChainLinkETHSource{}, dataSources[0])
	source0 := dataSources[0].(*chainlink.ChainLinkETHSource)
	require.Equal(t, "TOKEN-0", source0.Token)
	require.Equal(t, "UNIT-0", source0.Unit)
	require.Equal(t, uint8(8), source0.Decimals)

	require.IsType(t, &chainlink.ChainLinkETHSource{}, dataSources[1])
	source1 := dataSources[1].(*chainlink.ChainLinkETHSource)
	require.Equal(t, "TOKEN-1", source1.Token)
	require.Equal(t, "UNIT-1", source1.Unit)
	require.Equal(t, uint8(10), source1.Decimals)

	require.IsType(t, &RateLimitDataSource{}, dataSources[2])
	rateLimitSource2 := dataSources[2].(*RateLimitDataSource)
	require.IsType(t, &chainlink.ChainLinkETHSource{}, rateLimitSource2.Source)
	source2 := rateLimitSource2.Source.(*chainlink.ChainLinkETHSource)
	require.Equal(t, "TOKEN-2", source2.Token)
	require.Equal(t, "UNIT-2", source2.Unit)
	require.Equal(t, uint8(8), source2.Decimals)

	require.IsType(t, &RateLimitDataSource{}, dataSources[3])
	rateLimitSource3 := dataSources[3].(*RateLimitDataSource)
	require.IsType(t, &chainlink.ChainLinkETHSource{}, rateLimitSource3.Source)
	source3 := rateLimitSource3.Source.(*chainlink.ChainLinkETHSource)
	require.Equal(t, "TOKEN-3", source3.Token)
	require.Equal(t, "UNIT-3", source3.Unit)
	require.Equal(t, uint8(10), source3.Decimals)
}
