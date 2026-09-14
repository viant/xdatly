package tracing

import "net/http"

// SetStatusFromHTTPCode is the compatibility helper for the current
// transport-colored HTTP status mapping. The core tracing contract is the Trace
// and Span payload model plus explicit SetStatus(err); HTTP-code interpretation
// is not the governing center of the package.
func (s *Span) SetStatusFromHTTPCode(code int) {
	switch {
	case code >= 100 && code < 400:
		s.Status = SpanStatus{Code: StatusOK, Message: ""}
	case code >= 400 && code < 500:
		if s.Kind == "CLIENT" {
			s.Status = SpanStatus{Code: StatusError, Message: http.StatusText(code)}
			return
		}
		s.Status = SpanStatus{Code: StatusOK, Message: ""}
	case code >= 500:
		s.Status = SpanStatus{Code: StatusError, Message: http.StatusText(code)}
	default:
		s.Status = SpanStatus{Code: StatusOK, Message: ""}
	}
}
