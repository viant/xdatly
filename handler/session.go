package handler

import (
	"github.com/viant/xdatly/response"
)

// Session is the reduced per-invocation custom-handler facade.
// It stays intentionally small to avoid recreating the old fat session surface.
type Session interface {
	Binder() Binder
	Response() response.Writer
}
