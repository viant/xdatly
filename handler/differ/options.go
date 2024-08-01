package differ

type (
	//Options defins validator options
	Options struct {
		WithShallow   bool
		WithSetMarker bool
	}

	//Option defines validator option
	Option func(o *Options)
)

func (o *Options) Apply(opts ...Option) {
	if len(opts) == 0 {
		return
	}
	for _, opt := range opts {
		opt(o)
	}
}

// WithShallow creates with shallow option
func WithShallow(f bool) Option {
	return func(o *Options) {
		o.WithShallow = f
	}
}

// WithSetMarker creates with shallow option
func WithSetMarker(f bool) Option {
	return func(o *Options) {
		o.WithSetMarker = f
	}
}

type (
	logOptions struct {
		source string
		id     interface{}
		userID interface{}
	}

	LogOption func(l *logOptions)
)

func WithSource(source string) LogOption {
	return func(l *logOptions) {
		l.source = source
	}
}

func WithSourceID(id interface{}) LogOption {
	return func(l *logOptions) {
		l.id = id
	}
}

func WithUserID(userID interface{}) LogOption {
	return func(l *logOptions) {
		l.userID = userID
	}
}
