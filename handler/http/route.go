package http

// Route represents http route
type Route struct {
	URL    string
	Method string
}

// NewRoute creates a new route
func NewRoute(url, method string) *Route {
	return &Route{
		URL:    url,
		Method: method,
	}
}
