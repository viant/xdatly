package response

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/viant/xdatly/tracing"
)

// Logging/debug helpers over response metrics stay public for compatibility
// with current HTTP and logging surfaces, but are not the governing payload
// contract.

func (m *Metric) Name() string {
	return strings.Title(m.View)
}

func (m *Metric) HideSQL() *Metric {
	if m == nil {
		return nil
	}
	ret := *m
	if m.Executions == nil {
		return &ret
	}
	ret.Executions = make(SQLExecutions, len(m.Executions))
	copy(ret.Executions, m.Executions)
	for index, elem := range ret.Executions {
		if elem == nil {
			continue
		}
		copy := *elem
		copy.SQL = ""
		copy.Args = nil
		ret.Executions[index] = &copy
	}
	return &ret
}

func (m Metrics) HideMetrics() Metrics {
	result := make(Metrics, len(m))
	copy(result, m)
	for i, item := range m {
		result[i] = item.HideSQL()
	}
	return result
}

func (s *SQLExecution) ToSpan(viewName string) *tracing.Span {
	if s == nil {
		return nil
	}
	status := tracing.SpanStatus{Code: tracing.StatusOK}
	if s.Error != "" {
		status.Code = tracing.StatusError
		status.Message = strings.TrimSpace(s.Error)
	}

	attrs := map[string]string{
		"db.system": "sql",
		"db.rows":   strconv.Itoa(s.Rows),
	}
	if s.SQL != "" {
		attrs["db.statement"] = s.SQL
		attrs["db.args"] = fmt.Sprintf("%v", s.Args)
	}
	if s.CacheStats != nil {
		attrs["cache.key"] = fmt.Sprintf("%s", s.CacheStats.Key)
	}

	id := s.ID
	if id == "" {
		id = uuid.New().String()
	}
	name := "SQL Select: " + strings.Trim(viewName, "#")
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
		if metric == nil {
			continue
		}
		name := "Assemble: " + strings.Trim(metric.View, "#")
		switch metric.Type {
		case "INSERT":
			name = "SQL Insert: " + strings.Trim(metric.View, "#")
		case "UPDATE":
			name = "SQL Update: " + strings.Trim(metric.View, "#")
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
			if exec == nil {
				continue
			}
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
