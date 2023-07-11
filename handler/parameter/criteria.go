package parameter

import (
	"context"
	"reflect"
)

type (
	CriteriaBuilder interface {
		BuildCriteria(ctx context.Context, value interface{}, options *Options) error
	}

	Options struct {
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
		SetCriteria(query string, args []interface{})
		IgnoreRead()
	}
)
