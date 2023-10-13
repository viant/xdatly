package docs

import (
	"database/sql"
)

type (
	Options struct {
		URL       string
		Connector Connector
	}

	Option func(o *Options)

	Connector interface {
		DB() (*sql.DB, error)
	}
)

func WithURL(URL string) Option {
	return func(o *Options) {
		o.URL = URL
	}
}

func WithConnector(connector Connector) Option {
	return func(o *Options) {
		o.Connector = connector
	}
}
