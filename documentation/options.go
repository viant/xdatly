package documentation

type (
	Options struct {
		URL       string
		Connector string
	}

	Option func(o *Options)
)

func WithURL(URL string) Option {
	return func(o *Options) {
		o.URL = URL
	}
}

func WithConnector(connector string) Option {
	return func(o *Options) {
		o.Connector = connector
	}
}
