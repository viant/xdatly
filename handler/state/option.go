package state

import (
	"net/http"
	"net/url"
)

type (
	//Option represents state option
	Option func(o *Options)
	//Options represents state options
	Options struct {
		query      url.Values
		headers    http.Header
		body       []byte
		pathParams map[string]string

		scope       string
		form        *Form
		httpRequest *http.Request
		constants   map[string]interface{}
	}
)

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

// PathParameters returns path parameters
func (s *Options) PathParameters() map[string]string {
	return s.pathParams
}

// Query returns query
func (s *Options) Query() url.Values {
	return s.query
}

// Headers returns headers
func (s *Options) Headers() http.Header {
	return s.headers
}

// Body returns body
func (s *Options) Body() []byte {
	return s.body
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

// WithScope returns option with scope
func WithScope(scope string) Option {
	return func(o *Options) {
		o.scope = scope

	}
}

// WithPathParameter returns option with path parameters
func WithPathParameter(name, value string) Option {
	return func(o *Options) {
		if len(o.pathParams) == 0 {
			o.pathParams = make(map[string]string)
		}
		o.pathParams[name] = value
	}
}

// WithHttpRequest returns option with scope
func WithHttpRequest(httpRequest *http.Request) Option {
	return func(o *Options) {
		o.httpRequest = httpRequest

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

// WithQuery returns option with query
func WithQuery(query url.Values) Option {
	return func(o *Options) {
		o.query = query
	}
}

// WithQueryParameters returns option with query parameters
func WithQueryParameters(name string, values []string) Option {
	return func(o *Options) {
		if len(o.query) == 0 {
			o.query = make(url.Values)
		}
		o.query[name] = values
	}
}

// WithQueryParameter returns option with query parameters
func WithQueryParameter(name, value string) Option {
	return func(o *Options) {
		if len(o.query) == 0 {
			o.query = make(url.Values)
		}
		o.query[name] = []string{value}
	}
}

// WithHeader returns option with header
func WithHeader(name, value string) Option {
	return func(o *Options) {
		if len(o.headers) == 0 {
			o.headers = make(http.Header)
		}
		o.headers[name] = []string{value}
	}
}

// WithHeaders returns option with headers
func WithHeaders(headers http.Header) Option {
	return func(o *Options) {
		o.headers = headers
	}
}

// WithBody returns option with body
func WithBody(body []byte) Option {
	return func(o *Options) {
		o.body = body
	}
}
