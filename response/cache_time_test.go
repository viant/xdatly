package response

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCacheCreationTimeJSON(t *testing.T) {
	created := time.Date(2026, 9, 22, 12, 0, 0, 123000000, time.UTC)
	expiry := created.Add(time.Minute)
	data, err := json.Marshal(&CacheStats{CreatedTime: &created, ExpiryTime: &expiry})
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if fields["createdTime"] == nil {
		t.Fatal("createdTime missing")
	}
	var restored CacheStats
	if err = json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.CreatedTime == nil || !restored.CreatedTime.Equal(created) {
		t.Fatalf("creation time: %s", data)
	}
	if restored.ExpiryTime == nil || !restored.ExpiryTime.Equal(expiry) {
		t.Fatalf("expiry time: %s", data)
	}
	data, err = json.Marshal(&CacheStats{})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "{}" {
		t.Fatalf("unknown timestamps must be omitted: %s", data)
	}
}
