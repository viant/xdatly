package response

import "errors"

// StatusCoder is the minimal shared error/response contract for surfaces that
// carry an explicit status code.
type StatusCoder interface {
	StatusCode() int
}

// ErrorStatusCode unwraps err and returns its status code when present.
func ErrorStatusCode(err error, fallback int) int {
	if err == nil {
		return fallback
	}
	var coder StatusCoder
	if errors.As(err, &coder) {
		if code := coder.StatusCode(); code != 0 {
			return code
		}
	}
	return fallback
}
