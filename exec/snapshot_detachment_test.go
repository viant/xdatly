package exec

import (
	"sync"
	"testing"
	"time"

	"github.com/viant/xdatly/response"
	"github.com/viant/xdatly/tracing"
)

func TestCompletedSnapshotTiming(t *testing.T) {
	start := time.Now().Add(-time.Hour)
	for _, duration := range []time.Duration{0, 123 * time.Millisecond} {
		trace := tracing.NewTrace("svc", "v1")
		trace.Spans = []*tracing.Span{{StartTime: start}}
		c := New(WithStartTime(start), WithTrace(trace))
		c.Complete(start.Add(duration), nil)
		snapshot := c.SnapshotForLogging()
		if snapshot.ElapsedMs != int(duration.Milliseconds()) || snapshot.SnapshotForLogging().ElapsedMs != snapshot.ElapsedMs {
			t.Fatal("completed timing was recomputed")
		}
		if !snapshot.Trace.Spans[0].EndTime.Equal(start.Add(duration)) {
			t.Fatal("completed root span timing was lost")
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				c.Complete(start.Add(duration), nil)
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				if c.SnapshotForLogging().ElapsedMs != int(duration.Milliseconds()) {
					t.Error("completion timing changed")
				}
			}
		}()
		wg.Wait()
		c.Complete(start.Add(time.Second), nil)
		if snapshot.ElapsedMs != int(duration.Milliseconds()) || !snapshot.Trace.Spans[0].EndTime.Equal(start.Add(duration)) {
			t.Fatal("later completion changed frozen snapshot")
		}
	}
}

func TestSnapshotDetachesKnownPointers(t *testing.T) {
	parent := "parent"
	created := time.Unix(123, 0)
	expiry := time.Unix(456, 0)
	trace := tracing.NewTrace("svc", "v1")
	trace.Spans = []*tracing.Span{{ParentSpanID: &parent, Attributes: map[string]string{"key": "original"}}, nil}
	c := New(WithTrace(trace))
	c.AppendMetrics(&response.Metric{Executions: response.SQLExecutions{{Args: []any{"original"}, CacheStats: &response.CacheStats{CreatedTime: &created, ExpiryTime: &expiry}}, nil}})
	frozen := c.SnapshotForLogging()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			parent = "changed"
			created = time.Unix(int64(i), 0)
			expiry = created
			trace.Spans[0].Attributes["key"] = "changed"
			trace.Resource.ServiceName = "changed"
			c.Metrics[0].Executions[0].Args[0] = "changed"
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			exec := frozen.Metrics[0].Executions[0]
			if *frozen.Trace.Spans[0].ParentSpanID != "parent" || frozen.Trace.Spans[0].Attributes["key"] != "original" || frozen.Trace.Resource.ServiceName != "svc" || *exec.CacheStats.CreatedTime != time.Unix(123, 0) || *exec.CacheStats.ExpiryTime != time.Unix(456, 0) || exec.Args[0] != "original" {
				t.Error("snapshot retained source alias")
				return
			}
		}
	}()
	wg.Wait()
	if frozen.Trace.Spans[1] != nil || frozen.Metrics[0].Executions[1] != nil {
		t.Fatal("nil entries changed")
	}
}
