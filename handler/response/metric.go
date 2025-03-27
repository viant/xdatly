package response

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/viant/xdatly/handler/tracing"
	"strconv"
	"strings"
	"time"
)

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

func (m *Metric) ParametrizedSQL() []*ParametrizedSQL {
	var result = make([]*ParametrizedSQL, 0)
	if m.Executions == nil {
		return result
	}
	for _, tmpl := range m.Executions {
		result = append(result, &ParametrizedSQL{Query: tmpl.SQL, Args: tmpl.Args})
	}
	return result
}

func (m *Metric) Name() string {
	return strings.Title(m.View)
}

func (m *Metric) SQL() string {
	if m.Executions != nil && len(m.Executions) > 0 {
		tmpl := m.Executions[0]
		SQL := ExpandSQL(tmpl.SQL, tmpl.Args)
		return SQL
	}
	return ""
}

func (m *Metric) HideSQL() *Metric {
	ret := *m
	if m.Executions == nil {
		return &ret
	}
	ret.Executions = make(SQLExecutions, len(m.Executions))
	copy(ret.Executions, m.Executions)
	for _, elem := range m.Executions {
		elem.SQL = ""
		elem.Args = nil
	}
	return &ret
}

// Basic returns n
func (m Metrics) HideMetrics() Metrics {
	var result = make(Metrics, len(m))
	copy(result, m)
	for i, item := range m {
		m[i] = item.HideSQL()
	}
	return result
}

// SQL returns main view SQL
func (m Metrics) SQL() string {
	if m == nil || len(m) == 0 {
		return ""
	}
	return (m)[0].SQL()
}

// Lookup looks up view metric by name
func (m Metrics) Lookup(viewName string) *Metric {
	for _, candidate := range m {
		if candidate.View == viewName {
			return candidate
		}
	}
	return nil
}

func (m Metrics) ParametrizedSQL() []*ParametrizedSQL {
	var result []*ParametrizedSQL
	for _, metric := range m {
		result = append(result, metric.ParametrizedSQL()...)
	}
	return result
}

func (s *SQLExecution) ToSpan(viewName string) *tracing.Span {
	status := tracing.SpanStatus{
		Code: tracing.StatusOK,
	}
	if s.Error != "" {
		status.Code = tracing.StatusError
		status.Message = strings.TrimSpace(s.Error)
	}

	attrs := map[string]string{}
	attrs["db.system"] = "sql"
	attrs["db.rows"] = strconv.Itoa(s.Rows)
	if s.SQL != "" {
		attrs["db.statement"] = s.SQL
		attrs["db.args"] = fmt.Sprintf("%v", s.Args)
	}

	if s.CacheStats != nil {
		attrs["cache.key"] = fmt.Sprintf("%s", s.CacheStats.Key)
	}

	id := s.ID
	name := "SQL Select: " + strings.Trim(viewName, "#")
	if id == "" {
		id = uuid.New().String()
	}
	var parentID *string
	if s.ParentID != "" {
		parentID = &s.ParentID
	}
	return &tracing.Span{
		SpanID:       id,
		ParentSpanID: parentID,
		Name:         name,
		Kind:         "CLIENT",
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		Attributes:   attrs,
		Status:       status,
	}
}

func (m Metrics) ToSpans(ownerID *string) []*tracing.Span {
	var spans []*tracing.Span
	for _, metric := range m {
		name := ""
		switch metric.Type {
		case "INSERT":
			name = "SQL Insert: " + strings.Trim(metric.View, "#")
		case "UPDATE":
			name = "SQL Update: " + strings.Trim(metric.View, "#")
		default:
			name = "Assemble: " + strings.Trim(metric.View, "#")
		}
		span := tracing.Span{
			SpanID:       metric.ID,
			ParentSpanID: ownerID,
			Name:         name,
			Kind:         "CLIENT",
			StartTime:    metric.StartTime,
			EndTime:      metric.EndTime,
			Attributes: map[string]string{
				"elapsed":   metric.Elapsed,
				"elapsedMs": strconv.Itoa(metric.ElapsedMs),
				"rows":      strconv.Itoa(metric.Rows),
			},
			Status: tracing.SpanStatus{
				Code:    tracing.StatusOK,
				Message: metric.Error,
			},
		}
		if metric.Error != "" {
			span.Status.Code = tracing.StatusError
		}
		spans = append(spans, &span)
		for _, exec := range metric.Executions {
			execSpan := exec.ToSpan(metric.View)
			if execSpan.ParentSpanID == nil {
				execSpan.ParentSpanID = &metric.ID
			}
			if exec.Error != "" {
				execSpan.Status.Code = tracing.StatusError
			}
			if exec.CacheStats != nil {
				execSpan.Attributes["cacheType"] = exec.CacheStats.Type
				execSpan.Attributes["cacheKey"] = exec.CacheStats.Key
				execSpan.Attributes["cacheDataset"] = exec.CacheStats.Dataset
				execSpan.Attributes["cacheNamespace"] = exec.CacheStats.Namespace
				execSpan.Attributes["cacheRecordsCounter"] = strconv.Itoa(exec.CacheStats.RecordsCounter)
			}
			spans = append(spans, execSpan)
		}
	}
	return spans
}
