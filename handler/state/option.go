package state

import "net/http"

type (
	//Option represents state option
	Option func(o *Options)
	//Options represents state options
	Options struct {
		scope       string
		form        *Form
		httpRequest *http.Request
		constants   map[string]interface{}
	}
)

// WithScope returns option with scope
func WithScope(scope string) Option {
	return func(o *Options) {
		o.scope = scope

	}
}

// WithHttpRequest returns option with scope
func WithHttpRequest(httpRequest *http.Request) Option {
	return func(o *Options) {
		o.httpRequest = httpRequest

	}
}

// Scope returns scope
func (s *Options) Scope() string {
	return s.scope
}

// Form returns form
func (s *Options) Form() *Form {
	return s.form
}

// Constants returns constants
func (s *Options) Constants() map[string]interface{} {
	return s.constants
}

// HttpRequest returns http request
func (s *Options) HttpRequest() *http.Request {
	return s.httpRequest
}

// WithConstants returns option with constants
func WithConstants(key string, value string) Option {
	return func(o *Options) {
		if o.constants == nil {
			o.constants = make(map[string]interface{})
		}
		o.constants[key] = value
	}
}

// WithForm returns option with form
func WithForm(form *Form) Option {
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
