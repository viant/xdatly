package async

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/viant/xdatly/async/destination"
)

func TestJobJSONRoundTrip(t *testing.T) {
	errText := "failed"
	userEmail := "user@example.com"
	userID := "u-1"
	connector := "bq"
	tableName := "DATLY_JOBS"
	cacheName := "jobs"
	now := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC)
	start := now.Add(time.Second)
	end := start.Add(2 * time.Second)
	expiry := now.Add(time.Hour)
	create := destination.CreateDispositionIfNeeded
	write := destination.WriteDispositionAppend

	job := &Job{
		ID:       "job-1",
		MatchKey: "view:abc",
		Status:   StatusPending,
		Table: destination.Table{
			Connector:         &connector,
			TableName:         &tableName,
			CreateDisposition: &create,
			WriteDisposition:  &write,
		},
		Cache: destination.Cache{
			Cache: &cacheName,
		},
		Request: Request{
			Method: "GET",
			URI:    "/v1/api/jobs",
		},
		Principal: Principal{
			UserEmail: &userEmail,
			UserID:    &userID,
		},
		MainView:      "JobsView",
		Module:        "jobs",
		Labels:        "daily",
		JobType:       "read",
		EventURL:      "s3://bucket/key",
		Error:         &errText,
		CreationTime:  now,
		StartTime:     &start,
		EndTime:       &end,
		ExpiryTime:    &expiry,
		WaitTimeInMcs: 10,
		RunTimeInMcs:  20,
		Deactivated: true,
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var actual Job
	if err := json.Unmarshal(data, &actual); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if actual.Status != StatusPending {
		t.Fatalf("expected pending status, got %q", actual.Status)
	}
	if actual.Request.Method != "GET" || actual.Request.URI != "/v1/api/jobs" {
		t.Fatalf("unexpected request payload: %+v", actual.Request)
	}
	if actual.Connector == nil || *actual.Connector != "bq" {
		t.Fatalf("expected table connector to round-trip")
	}
	if actual.Cache.Cache == nil || *actual.Cache.Cache != "jobs" {
		t.Fatalf("expected cache destination to round-trip")
	}
}

func TestAsyncEnumsStayCompatible(t *testing.T) {
	if StatusPending != "PENDING" || StatusRunning != "RUNNING" || StatusDone != "DONE" || StatusError != "ERROR" {
		t.Fatalf("unexpected status constants")
	}
	if InvocationTypeEvent != "event" || InvocationTypeUndefined != "" {
		t.Fatalf("unexpected invocation type constants")
	}
	if NotificationMethodStorage != "Storage" || NotificationMethodMessageBus != "MessageBus" || NotificationMethodUndefined != "" {
		t.Fatalf("unexpected notification method constants")
	}
}

func TestJobCarriesNoSQLXPersistenceTags(t *testing.T) {
	jobType := reflect.TypeOf(Job{})
	for i := 0; i < jobType.NumField(); i++ {
		field := jobType.Field(i)
		if field.Tag.Get("sqlx") != "" || field.Tag.Get("sqlxAsync") != "" {
			t.Fatalf("unexpected persistence tag on %s", field.Name)
		}
	}
}
