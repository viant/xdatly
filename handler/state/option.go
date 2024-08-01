package state

import "github.com/viant/xdatly/handler/http"

type (
	//Option represents state option
	Option func(o *Options)
	//Options represents state options
	Options struct {
		scope     string
		form      *http.Form
		constants map[string]string
	}
)

// WithScope returns option with scope
func WithScope(scope string) Option {
	return func(o *Options) {
		o.scope = scope

	}
}

// Scope returns scope
func (s *Options) Scope() string {
	return s.scope
}

// Form returns form
func (s *Options) Form() *http.Form {
	return s.form
}

// Constants returns constants
func (s *Options) Constants() map[string]string {
	return s.constants
}

// WithConstants returns option with constants
func WithConstants(key string, value string) Option {
	return func(o *Options) {
		if o.constants == nil {
			o.constants = make(map[string]string)
		}
		o.constants[key] = value
	}
}

// WithForm returns option with form
func WithForm(form *http.Form) Option {
	return func(o *Options) {
		o.form = form
	}
}

// NewOptions returns option with scope
func NewOptions(options ...Option) *Options {
	var result = &Options{}
	for _, option := range options {
		option(result)
	}
	return result
}
