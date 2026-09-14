package docs

import "database/sql"

// Options is the reduced public docs-provider option bag.
type Options struct {
	URL       string
	Connector Connector
}

// Option mutates docs-provider options.
type Option func(o *Options)

// Connector is the reduced public docs connector contract.
//
// This is a bounded compatibility seam for provider-backed docs lookup, not a
// signal that docs ownership belongs to SQL/database runtime behavior.
type Connector interface {
	DB() (*sql.DB, error)
}

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
