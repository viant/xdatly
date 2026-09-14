package async

import (
	"time"

	"github.com/viant/xdatly/async/destination"
)

type Job struct {
	ID       string `json:"id,omitempty"`
	MatchKey string `json:"matchKey,omitempty"`
	Status   Status `json:"status,omitempty"`
	destination.Table
	destination.Cache
	Request
	Principal
	MainView      string     `json:"mainView,omitempty"`
	Module        string     `json:"module,omitempty"`
	Labels        string     `json:"labels,omitempty"`
	JobType       string     `json:"jobType,omitempty"`
	EventURL      string     `json:"eventUrl,omitempty"`
	Error         *string    `json:"error,omitempty"`
	CreationTime  time.Time  `json:"creationTime,omitempty"`
	StartTime     *time.Time `json:"startTime,omitempty"`
	EndTime       *time.Time `json:"endTime,omitempty"`
	ExpiryTime    *time.Time `json:"expiryTime,omitempty"`
	WaitTimeInMcs int        `json:"waitTimeInMcs,omitempty"`
	RunTimeInMcs  int        `json:"runTimeInMcs,omitempty"`
	Deactivated   bool       `json:"deactivated,omitempty"`
}

type Request struct {
	Method string `json:"method,omitempty"`
	URI    string `json:"uri,omitempty"`
}

type Principal struct {
	UserEmail *string `json:"userEmail,omitempty"`
	UserID    *string `json:"userId,omitempty"`
}
