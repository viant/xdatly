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

	JobInfo struct {
		RequestTime time.Time `json:",omitempty"`
		JobStatus   string    `json:",omitempty"`
		CreateTime  time.Time `json:",omitempty"`
		WaitTimeMcs int       `json:",omitempty"`
		RunTimeMcs  int       `json:",omitempty"`
		ExpiryInSec int       `json:",omitempty"`
		CacheKey    string    `json:",omitempty"`
		CacheHit    bool      `json:",omitempty"`
		MatchKey    string    `json:",omitempty"`
	}
)
