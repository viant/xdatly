package async

import (
	"time"
)

type (
	Async interface {
		Read(options ...ReadOption) (*Record, error)
		ReadInto(dst interface{}, record *DbRecord) error
	}

	ReadOptions struct {
		Connector   string
		Record      *DbRecord
		OnExist     *OnExist
		Destination string
	}

	OnExist struct {
		Return  bool
		Refresh time.Duration
	}

	ReadOption func(options *ReadOptions)
)

func WithReadOptions(options ReadOptions) ReadOption {
	return func(opts *ReadOptions) {
		*opts = options
	}
}

func WithConnector(name string) ReadOption {
	return func(options *ReadOptions) {
		options.Connector = name
	}
}

func WithRecord(record *DbRecord) ReadOption {
	return func(options *ReadOptions) {
		options.Record = record
	}
}

func WithOnExist(onExist *OnExist) ReadOption {
	return func(options *ReadOptions) {
		options.OnExist = onExist
	}
}

func WithDestination(URL string) ReadOption {
	return func(options *ReadOptions) {
		options.Destination = URL
	}
}
