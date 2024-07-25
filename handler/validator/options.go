package validator

import "database/sql"

type CanUseMarkerProvider func(v interface{}) bool

type (
	//Options defins validator options
	Options struct {
		WithShallow          bool
		WithSetMarker        bool
		WithDB               *sql.DB
		WithUnique           bool
		WithRef              bool
		Location             string
		CanUseMarkerProvider CanUseMarkerProvider
		WithValidation       *Validation
	}

	//Option defines validator option
	Option func(o *Options)
)

func (o *Options) Apply(opts []Option) {
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

// WithShallow creates with shallow option
func WithSetMarker(f bool) Option {
	return func(o *Options) {
		o.WithSetMarker = f
	}
}

// WithDB create with db validation option
func WithDB(db *sql.DB) Option {
	return func(o *Options) {
		o.WithDB = db
	}
}

// WithRefCheck return with ref check option
func WithRefCheck(flag bool) Option {
	return func(o *Options) {
		o.WithRef = flag
	}
}

// WithUnique returns with unique check
func WithUnique(flag bool) Option {
	return func(o *Options) {
		o.WithUnique = flag
	}
}

// WithValidation creates with Validation option
func WithValidation(v *Validation) Option {
	return func(o *Options) {
		o.WithValidation = v
	}
}

// WithLocation creates with location option
func WithLocation(location string) Option {
	return func(o *Options) {
		o.Location = location
	}
}

// WithCanUseMarkerProvider creates with marker provider option
func WithCanUseMarkerProvider(provider CanUseMarkerProvider) Option {
	return func(o *Options) {
		o.WithSetMarker = true
		o.CanUseMarkerProvider = provider
	}
}
