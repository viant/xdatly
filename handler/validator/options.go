package validator

type (
	//Options defins validator options
	Options struct {
		WithShallow    bool
		WithSetMarker  bool
		WithValidation *Validation
	}

	//Option defines validator option
	Option func(o *Options)
)

//WithShallow creates with shallow option
func WithShallow(f bool) Option {
	return func(o *Options) {
		o.WithShallow = f
	}
}

//WithShallow creates with shallow option
func WithSetMarker(f bool) Option {
	return func(o *Options) {
		o.WithSetMarker = f
	}
}

//WithValidation creates with Validation option
func WithValidation(v *Validation) Option {
	return func(o *Options) {
		o.WithValidation = v
	}
}
