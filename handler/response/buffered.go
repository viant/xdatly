package response

import (
	"bytes"
	"io"
	"net/http"
)

// Buffered represents a buffered response
type Buffered struct {
	buffer      *bytes.Buffer
	statusCode  int
	compression string
	size        int
	headers     http.Header
}

func (b *Buffered) StatusCode() int {
	return b.statusCode
}

func (b *Buffered) Body() io.Reader {
	return b.buffer
}

func (b *Buffered) Headers() http.Header {
	return b.headers
}

func (b *Buffered) Size() int {
	return b.size
}

func (b *Buffered) CompressionType() string {
	return b.compression
}

func (b *Buffered) SetStatusCode(status int) {
	b.statusCode = status
}

func NewBuffered(options ...Option) *Buffered {
	response := &Buffered{}
	for _, option := range options {
		option(response)
	}
	if response.buffer == nil {
		response.buffer = &bytes.Buffer{}
	}
	if response.size == 0 {
		response.size = response.buffer.Len()
	}
	return response
}

// Option represents a buffered response option
type Option func(r *Buffered)

// Options represents a list of buffered response options
type Options []Option

// AdjustStatusCode adjusts the status code
func (o *Options) AdjustStatusCode(candidates ...interface{}) {
	for _, candidate := range candidates {
		if candidate == nil {
			continue
		}
		if coder, ok := candidate.(StatusCoder); ok {
			*o = append(*o, WithStatusCode(coder.StatusCode()))
		}
	}
}

// Append appends options
func (o *Options) Append(opts ...Option) {
	*o = append(*o, opts...)
}

func (o *Options) Options() Options {
	return *o
}

// WithStatusCode returns an option that sets the status code
func WithStatusCode(statusCode int) Option {
	return func(r *Buffered) {
		r.statusCode = statusCode
	}
}

// WithBytes returns an option that sets the payload
func WithBytes(bs []byte) Option {
	return func(r *Buffered) {
		r.buffer = bytes.NewBuffer(bs)
	}
}

// WithBuffer returns an option that sets the buffer
func WithBuffer(buffer *bytes.Buffer) Option {
	return func(r *Buffered) {
		r.buffer = buffer
	}
}

// WithCompressions returns an option that sets the compression type
func WithCompressions(compressions string) Option {
	return func(r *Buffered) {
		r.compression = compressions
	}
}

// WithSize returns an option that sets the size
func WithSize(compressions string) Option {
	return func(r *Buffered) {
		r.compression = compressions
	}
}

// WithHeader returns an option that sets the header
func WithHeader(name, value string) Option {
	return func(r *Buffered) {
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		r.headers.Add(name, value)
	}
}

// WithHeaders returns an option that sets the headers
func WithHeaders(header http.Header) Option {
	return func(r *Buffered) {
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		for key, values := range header {
			header.Set(key, values[0])
		}
	}
}
