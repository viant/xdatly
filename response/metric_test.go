package response

import (
	"testing"
	"time"
)

func TestMetricContracts(t *testing.T) {
	now := time.Now()
	metrics := Metrics{
		&Metric{
			ID:        "metric-1",
			StartTime: now,
			EndTime:   now,
			View:      "users",
			Type:      "SELECT",
			Elapsed:   "1ms",
			ElapsedMs: 1,
			Rows:      2,
			Executions: SQLExecutions{
				&SQLExecution{
					ID:        "sql-1",
					StartTime: now,
					EndTime:   now,
					SQL:       "SELECT * FROM users WHERE id = ?",
					Args:      []interface{}{7},
					Rows:      2,
					CacheStats: &CacheStats{
						Type:           "aerospike",
						Key:            "users:7",
						RecordsCounter: 2,
					},
				},
			},
		},
	}

	if metrics.SQL() == "" {
		t.Fatalf("expected SQL expansion")
	}
	if metrics.Lookup("users") == nil {
		t.Fatalf("expected view lookup")
	}
	if len(metrics.ParametrizedSQL()) != 1 {
		t.Fatalf("expected one parametrized sql")
	}
	if len(metrics.ToSpans(nil)) != 2 {
		t.Fatalf("expected metric and sql execution spans")
	}
	hidden := metrics.HideMetrics()
	if hidden[0].Executions[0].SQL != "" {
		t.Fatalf("expected hidden SQL")
	}
}
