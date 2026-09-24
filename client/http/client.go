// Package http defines public outbound HTTP capabilities without runtime or
// concrete transport dependencies. Clients are borrowed from injected providers.
package http

import "net/http"

// Client executes explicit HTTP requests. Request contexts control cancellation.
// The caller closes each response body; the provider owner manages connections.
type Client interface {
	Do(request *http.Request) (*http.Response, error)
}
