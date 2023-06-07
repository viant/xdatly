package sqlx

import "database/sql"

type Options struct {
	WithConnector string
	WithDb        *sql.DB
	WithTx        *sql.Tx
}

type Option func(o *Options)

func WithConnector(name string) Option {
	return func(o *Options) {
		o.WithConnector = name
	}
}

func WithSqlDB(db *sql.DB) Option {
	return func(o *Options) {
		o.WithDb = db
	}
}
func WithSqlTx(tx *sql.Tx) Option {
	return func(o *Options) {
		o.WithTx = tx
	}
}
