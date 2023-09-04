package response

type (
	Warning struct {
		Message string `json:",omitempty"`
		Reason  string `json:",omitempty"`
	}
	ResponseStatus struct {
		Status  string                 `json:",omitempty"`
		Message string                 `json:",omitempty"`
		Errors  interface{}            `json:",omitempty"`
		Warning []*Warning             `json:",omitempty"`
		Extras  map[string]interface{} `json:",omitempty" default:"embedded=true"`
	}
)
