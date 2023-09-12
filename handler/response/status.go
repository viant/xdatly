package response

import "time"

type (
	Warning struct {
		Message string `json:",omitempty"`
		Reason  string `json:",omitempty"`
	}
	Status struct {
		Status  string                 `json:",omitempty"`
		Message string                 `json:",omitempty"`
		Errors  interface{}            `json:",omitempty"`
		Warning []*Warning             `json:",omitempty"`
		Extras  map[string]interface{} `json:",omitempty" default:"embedded=true"`
	}

	JobStatus struct {
		RequestTime time.Time
		JobStatus   string
		CreateTime  time.Time
		WaitTimeMcs int
		RunTimeMcs  int
		ExpiryInSec int
		CacheKey    string
		CacheHit    bool
	}
)
