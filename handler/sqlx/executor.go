package sqlx

type Executor interface {
	Flusher
	Execute(DML string) error
}

type ExecutorOptions struct {
	Args []interface{}
}

type ExecutorOption func(o *ExecutorOptions)

func WithExecutorArgs(args ...interface{}) ExecutorOption {
	return func(o *ExecutorOptions) {
		o.Args = args
	}
}
