package connector

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestProviderRemainsOneExplicitNamedCapability(t *testing.T) {
	contract := reflect.TypeFor[Provider]()
	if contract.NumMethod() != 1 {
		t.Fatalf("connector contract has %d methods", contract.NumMethod())
	}
	method := contract.Method(0)
	if method.Name != "Connector" || method.Type.NumIn() != 2 || method.Type.In(0) != reflect.TypeFor[context.Context]() || method.Type.In(1) != reflect.TypeFor[string]() || method.Type.NumOut() != 2 || method.Type.Out(0) != reflect.TypeFor[*sql.DB]() || method.Type.Out(1) != reflect.TypeFor[error]() {
		t.Fatalf("unexpected connector contract: %s %v", method.Name, method.Type)
	}
}
