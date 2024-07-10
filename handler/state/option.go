package state

type (
	//Option represents state option
	Option func(o *Options)
	//Options represents state options
	Options struct {
		scope string
	}
)

// Scope returns scope
func (s *Options) Scope() string {
	return s.scope
}

// NewOptions returns option with scope
func NewOptions(options ...Option) *Options {
	var result = &Options{}
	for _, option := range options {
		option(result)
	}
	return result
}
