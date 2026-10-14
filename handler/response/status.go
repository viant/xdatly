package response

// Warning is the original data-only response warning contract.
type Warning struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// Status is the original data-only full response envelope.
type Status struct {
	Status  string                 `json:"status,omitempty"`
	Message string                 `json:"message,omitempty"`
	Errors  interface{}            `json:"errors,omitempty"`
	Warning []*Warning             `json:"warning,omitempty"`
	Extras  map[string]interface{} `json:"extras,omitempty" format:"inline=true"`
}
