package handler

import "testing"

func TestCapabilityKeys(t *testing.T) {
	if LoggerKey != "logger" || ValidatorKey != "validator" || MessageBusKey != "mbus" || ConnectorKey != "connector" {
		t.Fatalf("unexpected capability keys: %q %q %q %q", LoggerKey, ValidatorKey, MessageBusKey, ConnectorKey)
	}
}
