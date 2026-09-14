package response

// Status is the minimal shared output-status payload contract.
// It stays intentionally small so runtimes can populate it without pulling
// response-shaping internals into the public surface.
type Status struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   any    `json:"error,omitempty"`
}
