package http

// Route represents http route
type Route struct {
	Method string
	URL    string
}

// NewRoute creates a new route
func NewRoute(method, url string) *Route {
	return &Route{
		Method: method,
		URL:    url,
	}
}
