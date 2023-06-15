package sqlx

import (
	"context"
	"database/sql"
	"github.com/viant/xdatly/handler/validator"
)

type Sqlx interface {
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

type Service struct {
	sqlx Sqlx
}

func (s Service) Load(tableName string, data interface{}) error {
	return s.sqlx.Load(tableName, data)
}

func (s Service) Allocate(ctx context.Context, tableName string, dest interface{}, selector string) error {
	return s.Allocate(ctx, tableName, dest, selector)
}

func (s Service) Flush(ctx context.Context, tableName string) error {
	return s.sqlx.Flush(ctx, tableName)
}

func (s Service) Insert(tableName string, data interface{}) error {
	return s.sqlx.Insert(tableName, data)
}

func (s Service) Update(tableName string, data interface{}) error {
	return s.sqlx.Update(tableName, data)
}

func (s Service) Delete(tableName string, data interface{}) error {
	return s.sqlx.Delete(tableName, data)
}

func (s Service) Execute(DML string, options ...ExecutorOption) error {
	return s.sqlx.Execute(DML, options...)
}

func (s Service) Read(ctx context.Context, dest interface{}, SQL string, params ...interface{}) error {
	return s.sqlx.Read(ctx, dest, SQL, params)
}

func (s Service) Db(ctx context.Context) (*sql.DB, error) {
	return s.sqlx.Db(ctx)
}

func (s Service) Tx(ctx context.Context) (*sql.Tx, error) {
	return s.sqlx.Tx(ctx)
}

func (s Service) Validator() validator.Service {
	return s.sqlx.Validator()
}

func New(sqlx Sqlx) *Service {
	return &Service{sqlx: sqlx}
}
