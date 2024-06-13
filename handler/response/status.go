package response

import "time"

type (
	Warning struct {
		Message string `json:"message,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	Status struct {
		Status  string                 `json:",omitempty"`
		Message string                 `json:",omitempty"`
		Errors  interface{}            `json:",omitempty"`
		Warning []*Warning             `json:",omitempty"`
		Extras  map[string]interface{} `json:",omitempty" format:"inline=true"`
	}

	JobInfo struct {
		RequestTime  time.Time `json:",omitempty"`
		JobStatus    string    `json:",omitempty"`
		CreateTime   time.Time `json:",omitempty"`
		WaitTimeInMs int       `json:",omitempty"`
		RunTimeInMs  int       `json:",omitempty"`
		ExpiryInSec  int       `json:",omitempty"`
		CacheKey     string    `json:",omitempty"`
		CacheHit     bool      `json:",omitempty"`
		MatchKey     string    `json:",omitempty"`
	}
)
