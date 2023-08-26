package response

import "github.com/viant/xdatly/handler/validator"

type (
	Warning struct {
		Message string `json:",omitempty"`
		Reason  string `json:",omitempty"`
	}
	Status struct {
		Status  string                 `json:",omitempty"`
		Message string                 `json:",omitempty"`
		Errors  []*validator.Violation `json:",omitempty"`
		Warning []*Warning             `json:",omitempty"`
		Extras  map[string]interface{} `json:",omitempty" default:"embedded=true"`
	}
)
