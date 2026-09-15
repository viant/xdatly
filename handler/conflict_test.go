package handler_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/viant/xdatly/handler"
	"github.com/viant/xdatly/response"
)

func TestConflictPublicContract(t *testing.T) {
	cause := &handler.Conflict{Entity: "Orders", Field: "Version", Reason: "expected token is missing"}
	err := fmt.Errorf("validate: %w", cause)
	var conflict *handler.Conflict
	if !errors.As(err, &conflict) || conflict != cause || response.ErrorStatusCode(err, 500) != 409 {
		t.Fatal("conflict type or status lost", err)
	}
	if handler.WriteDelete != handler.WriteAction("delete") {
		t.Fatal("delete write action changed")
	}
}
