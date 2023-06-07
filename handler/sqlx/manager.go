package sqlx

import (
	"context"
	"database/sql"
	"github.com/viant/xdatly/handler/validator"
)

type Service interface {
	Sequencer
	Inserter
	Loader
	Updater
	Deleter
	Executor
	Reader
	Db(ctx context.Context) (*sql.DB, error)
	Tx(ctx context.Context) (*sql.Tx, error)
	Validator() validator.Service
}
