package response

import "time"

type (
	Metric struct {
		ID         string        `json:"id,omitempty"`
		StartTime  time.Time     `json:"startTime,omitempty"`
		EndTime    time.Time     `json:"endTime,omitempty"`
		View       string        `json:"view,omitempty"`
		Type       string        `json:"type,omitempty"`
		Elapsed    string        `json:"elapsed,omitempty"`
		ElapsedMs  int           `json:"elapsedMs,omitempty"`
		Rows       int           `json:"rows,omitempty"`
		Executions SQLExecutions `json:"executions,omitempty"`
		Error      string        `json:"error,omitempty"`
	}

	Metrics    []*Metric
	CacheStats struct {
		Type           string     `json:"type,omitempty"`
		RecordsCounter int        `json:"recordsCounter,omitempty"`
		Key            string     `json:"key,omitempty"`
		Dataset        string     `json:"dataset,omitempty"`
		Namespace      string     `json:"namespace,omitempty"`
		FoundWarmup    bool       `json:"foundWarmup,omitempty"`
		FoundLazy      bool       `json:"foundLazy,omitempty"`
		ErrorType      string     `json:"errorType,omitempty"`
		ErrorCode      int        `json:"errorCode,omitempty"`
		ExpiryTime     *time.Time `json:"expiryTime,omitempty"`
		CreatedTime    *time.Time `json:"createdTime,omitempty"`
	}

	SQLExecutions []*SQLExecution

	SQLExecution struct {
		ID         string        `json:"id,omitempty"`
		ParentID   string        `json:"parentID,omitempty"`
		StartTime  time.Time     `json:"startTime,omitempty"`
		EndTime    time.Time     `json:"endTime,omitempty"`
		SQL        string        `json:"sql,omitempty"`
		Args       []interface{} `json:"args,omitempty"`
		Rows       int           `json:"rows,omitempty"`
		CacheStats *CacheStats   `json:"cacheStats,omitempty"`
		Error      string        `json:"error,omitempty"`
	}

	ParametrizedSQL struct {
		Query string
		Args  []interface{}
	}
)

func (s *SQLExecution) SetError(err error) {
	if err == nil {
		return
	}
	s.Error = err.Error()
}

func (s *SQLExecutions) Append(executions ...*SQLExecution) {
	*s = append(*s, executions...)
}

func (m *Metrics) Append(metrics ...*Metric) {
	*m = append(*m, metrics...)
}

func (m Metrics) Lookup(viewName string) *Metric {
	for _, candidate := range m {
		if candidate.View == viewName {
			return candidate
		}
	}
	return nil
}

// SetCreatedTime records a detached timestamp from a native cache entry.
func (s *CacheStats) SetCreatedTime(value *time.Time) {
	s.CreatedTime = nil
	if value != nil {
		copy := *value
		s.CreatedTime = &copy
	}
}
