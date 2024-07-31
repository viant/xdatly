package response

import (
	"io"
	"net/http"
)

// StatusCoder represents a status code
type StatusCoder interface {
	StatusCode() int
	SetStatusCode(int)
}

// Response represents a response
type Response interface {
	StatusCoder
	Body() io.Reader
	Headers() http.Header
	Size() int
}

// Compressed represents a compressed response
type Compressed interface {
	CompressionType() string
}
