package handler

import (
	"context"
	"fmt"
	"reflect"
)

// RowConsumer receives each decoded row of one explicitly selected reader.
// It is invocation-local and is never inherited by dependency or child reads.
// The reader does not retain rows in its output. A callback must not perform
// nested database work while the source cursor is open. Its error stops the
// read; already delivered callbacks are never replayed by a database retry.
type RowConsumer struct {
	rowType reflect.Type
	consume func(context.Context, any) error
}

// ConsumeRows creates a typed consumer; T must be the generated row pointer.
// Keeping a row in the callback is the caller's explicit memory ownership.
func ConsumeRows[T any](consume func(context.Context, T) error) *RowConsumer {
	result := &RowConsumer{rowType: reflect.TypeFor[T]()}
	if consume != nil {
		result.consume = func(ctx context.Context, value any) error {
			row, ok := value.(T)
			if !ok {
				return fmt.Errorf("row consumer expects %v, got %T", result.rowType, value)
			}
			return consume(ctx, row)
		}
	}
	return result
}
func (c *RowConsumer) RowType() reflect.Type {
	if c == nil {
		return nil
	}
	return c.rowType
}
func (c *RowConsumer) Consume(ctx context.Context, row any) error {
	if c == nil || c.consume == nil {
		return fmt.Errorf("row consumer callback is required")
	}
	return c.consume(ctx, row)
}
func (c *RowConsumer) Valid() bool { return c != nil && c.rowType != nil && c.consume != nil }
