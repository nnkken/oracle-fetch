package datasource

import (
	"fmt"

	"go.uber.org/ratelimit"

	"github.com/nnkken/oracle-fetch/datasource/chainlink-eth"
	"github.com/nnkken/oracle-fetch/datasource/types"
)

var dataSourceFactoryMap = map[string]types.DataSourceFactory{
	"chainlink-eth": chainlink.NewDataSourceFromConfig,
}

func constructDecorator(config types.DataSourceConfig) types.DataSourceDecorator {
	var decorators []types.DataSourceDecorator
	if config.RateLimitRps > 0 {
		limiter := ratelimit.New(config.RateLimitRps)
		decorators = append(decorators, RateLimitDecorator(limiter))
	}
	return func(source types.DataSource) types.DataSource {
		for _, decorator := range decorators {
			source = decorator(source)
		}
		return source
	}
}

func InitDataSourcesFromConfig(config []types.DataSourceConfig) ([]types.DataSource, error) {
	var dataSources []types.DataSource
	for _, c := range config {
		factory, ok := dataSourceFactoryMap[c.Type]
		if !ok {
			return nil, fmt.Errorf("unknown data source type: %s", c.Type)
		}
		sources, err := factory(c.Config)
		if err != nil {
			return nil, fmt.Errorf("error when initializing data source for type %s: %w", c.Type, err)
		}
		decorate := constructDecorator(c)
		for _, source := range sources {
			dataSources = append(dataSources, decorate(source))
		}
	}
	return dataSources, nil
}
