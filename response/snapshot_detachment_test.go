package response

import (
	"sync"
	"testing"
	"time"
)

func TestLoggingProjectionsDetachKnownPointers(t *testing.T) {
	created, expiry := time.Unix(123, 0), time.Unix(456, 0)
	owner := "owner"
	metrics := Metrics{{ID: "metric", Executions: SQLExecutions{{ParentID: "explicit", CacheStats: &CacheStats{CreatedTime: &created, ExpiryTime: &expiry}}, {}}}}
	hidden := metrics.HideMetrics()
	spans := metrics.ToSpans(&owner)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			owner = "changed"
			metrics[0].ID = "changed"
			metrics[0].Executions[0].ParentID = "changed"
			metrics[0].Executions[0].CacheStats.Key = "changed"
			created = time.Unix(int64(i), 0)
			expiry = created
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			cache := hidden[0].Executions[0].CacheStats
			if *spans[0].ParentSpanID != "owner" || *spans[1].ParentSpanID != "explicit" || *spans[2].ParentSpanID != "metric" || cache.Key != "" || *cache.CreatedTime != time.Unix(123, 0) || *cache.ExpiryTime != time.Unix(456, 0) {
				t.Error("projection retained source alias")
				return
			}
		}
	}()
	wg.Wait()
}
