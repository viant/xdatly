package response

import (
	"io"
	"net/http"
)

// Response is a transport-ready response that an HTTP surface can write
// directly without JSON shaping.
type Response interface {
	StatusCoder
	Body() io.Reader
	Headers() http.Header
	Size() int
	SetStatusCode(int)
}

// Compressed marks a transport response whose body is already encoded.
type Compressed interface {
	CompressionType() string
}
