package codec

import (
	"context"
	"reflect"
)

type (
	criteriaBuilderKey string

	CriteriaBuilder interface {
		BuildCriteria(ctx context.Context, value interface{}, options *CriteriaBuilderOptions) (*Criteria, error)
	}

	Criteria struct {
		Expression   string
		Placeholders []interface{}
	}

	CriteriaBuilderOptions struct {
		Expression string
		Columns    ColumnsSource
		Parameters ValueGetter
		Selector   Selector
	}

	ColumnsSource interface {
		Column(key string) (Column, bool)
		ColumnName(key string) (string, error)
	}

	Column interface {
		ColumnName() string
		ColumnType() reflect.Type
		FieldName() string
	}

	ValueGetter interface {
		Value(ctx context.Context, paramName string) (interface{}, error)
	}

	Selector interface {
		IgnoreRead()
	}
)

const (
	_criteriaBuilderKey criteriaBuilderKey = "CriteriaBuilder"
)

// NewCriteriaBuilder creates a new criteria builder context
func NewCriteriaBuilder(ctx context.Context, builder CriteriaBuilder) context.Context {
	return context.WithValue(ctx, _criteriaBuilderKey, builder)
}

// CriteriaBuilderFromContext returns criteria builder from context
func CriteriaBuilderFromContext(ctx context.Context) CriteriaBuilder {
	if ctx == nil {
		return nil
	}
	if builder, ok := ctx.Value(_criteriaBuilderKey).(CriteriaBuilder); ok {
		return builder
	}
	return nil
}
