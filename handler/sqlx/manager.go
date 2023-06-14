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

func (s Service) Allocate(ctx context.Context, tableName string, dest interface{}, selector string) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Flush(ctx context.Context, tableName string) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Insert(tableName string, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Update(tableName string, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Delete(tableName string, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Execute(DML string, options ...ExecutorOption) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Read(ctx context.Context, dest interface{}, SQL string, params ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Db(ctx context.Context) (*sql.DB, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Tx(ctx context.Context) (*sql.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Validator() validator.Service {
	//TODO implement me
	panic("implement me")
}

func New(sqlx Sqlx) *Service {
	return &Service{sqlx: sqlx}
}
