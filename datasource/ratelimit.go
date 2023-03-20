package datasource

import (
	"go.uber.org/ratelimit"

	"github.com/nnkken/oracle-fetch/types"
)

var _ types.DataSource = (*RateLimitDataSource)(nil)

type RateLimitDataSource struct {
	source  types.DataSource
	limiter ratelimit.Limiter
}

func NewRateLimitDataSource(source types.DataSource, limiter ratelimit.Limiter) *RateLimitDataSource {
	return &RateLimitDataSource{
		source:  source,
		limiter: limiter,
	}
}

func (s *RateLimitDataSource) Fetch() ([]types.DBEntry, error) {
	s.limiter.Take()
	return s.source.Fetch()
}
