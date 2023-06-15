package response

import "net/http"

type Response interface {
	StatusCode() int
	Headers() http.Header
	Value() interface{}
}
