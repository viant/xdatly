package logger

import (
	"context"
	"testing"
)

func TestInvocationContextLoggerIdentity(t *testing.T) {
	var first, second Logger = &contextTestLogger{}, &contextTestLogger{}
	base := context.Background()
	parent := WithContext(base, first)
	child := WithContext(parent, second)
	if FromContext(nil) != nil || FromContext(base) != nil || FromContext(parent) != first || FromContext(child) != second || FromContext(WithContext(child, nil)) != nil {
		t.Fatal("logger context identity or shadowing failed")
	}
}

type contextTestLogger struct{ marker int }

func (*contextTestLogger) Debug(string, ...any) {}
func (*contextTestLogger) Info(string, ...any)  {}
func (*contextTestLogger) Warn(string, ...any)  {}
func (*contextTestLogger) Error(string, ...any) {}
