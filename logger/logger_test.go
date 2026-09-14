package logger

import "testing"

type testLogger struct{}

func (testLogger) Debug(string, ...any) {}
func (testLogger) Info(string, ...any)  {}
func (testLogger) Warn(string, ...any)  {}
func (testLogger) Error(string, ...any) {}

func TestLoggerContractCompile(t *testing.T) {
	var _ Logger = testLogger{}
}
