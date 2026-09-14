package response

import "testing"

func TestStatusPayloadShape(t *testing.T) {
	status := Status{
		Status:  "error",
		Message: "boom",
		Error:   "boom",
	}

	if status.Status != "error" {
		t.Fatalf("expected status field to be preserved")
	}
	if status.Message != "boom" {
		t.Fatalf("expected message field to be preserved")
	}
	if status.Error != "boom" {
		t.Fatalf("expected error field to be preserved")
	}
}
