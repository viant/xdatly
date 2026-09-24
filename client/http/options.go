package http

import "github.com/viant/xdatly/client"

// Options describes an explicit HTTP destination and its defaults. Dynamic
// credentials belong on the request, not in deployment constants.
type Options struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header,omitempty"`
	client.Limits
}
