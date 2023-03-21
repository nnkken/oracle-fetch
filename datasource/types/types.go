package types

import (
	"encoding/json"

	"github.com/nnkken/oracle-fetch/db"
)

type DataSource interface {
	Fetch() ([]db.DBEntry, error)
}

type DataSourceConfig struct {
	Type         string          `json:"type"`
	RateLimitRps int             `json:"rate_limit_rps"`
	Config       json.RawMessage `json:"config"`
}

type DataSourceDecorator = func(DataSource) DataSource

type DataSourceFactory = func(json.RawMessage) ([]DataSource, error)
