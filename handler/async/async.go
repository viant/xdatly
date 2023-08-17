package async

import (
	"context"
	"time"
)

const (
	HandlerTypeS3        HandlerType = "S3"
	HandlerTypeSQS       HandlerType = "SQS"
	HandlerTypeUndefined HandlerType = ""
)

type (
	HandlerType string
	Async       interface {
		Read(ctx context.Context, config *Config, options ...ReadOption) (*JobWithMeta, error)
		ReadInto(ctx context.Context, dst interface{}, job *Job, connector string) error
	}

	ReadOptions struct {
		Connector string
		Job       *Job
		OnExist   *OnExist
	}

	OnExist struct {
		Return  bool
		Refresh time.Duration
	}

	Config struct {
		HandlerType HandlerType
		BucketURL   string //S3
		QueueName   string //SQS
		Dataset     string //Bigquery db
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

func WithRecord(record *Job) ReadOption {
	return func(options *ReadOptions) {
		options.Job = record
	}
}

func WithOnExist(onExist *OnExist) ReadOption {
	return func(options *ReadOptions) {
		options.OnExist = onExist
	}
}
