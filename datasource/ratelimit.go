package datasource

import (
	"go.uber.org/ratelimit"

	"github.com/nnkken/oracle-fetch/types"
)

var _ types.DataSource = (*RateLimitDataSource)(nil)

type RateLimitDataSource struct {
	Source  types.DataSource
	Limiter ratelimit.Limiter
}

func RateLimitDecorator(limiter ratelimit.Limiter) types.DataSourceDecorator {
	return func(source types.DataSource) types.DataSource {
		return &RateLimitDataSource{
			Source:  source,
			Limiter: limiter,
		}
	}
}

func (s *RateLimitDataSource) Fetch() ([]types.DBEntry, error) {
	s.Limiter.Take()
	return s.Source.Fetch()
}
