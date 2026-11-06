package contract_test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/viant/bindly"
	"github.com/viant/xdatly/contract"
)

type sparseRecord struct {
	ID           string  `json:"id" sqlx:"ID,primaryKey" required:"true"`
	PricingModel int     `json:"pricingModel" required:"true" choice:"0:fixed" choice:"10:usage"`
	Commission   float64 `json:"commission" required:"true"`
}

func (*sparseRecord) Contract(ctx context.Context, binder bindly.Binder) (*contract.Policy, error) {
	var operation struct {
		Method string `parameter:",kind=contract,in=method,required=true"`
	}
	if err := binder.Bind(ctx, &operation); err != nil {
		return nil, err
	}
	if operation.Method != http.MethodPatch {
		return nil, nil
	}
	return &contract.Policy{Name: "record:patch:v1", Field: func(_ string, _ reflect.Type, field *contract.Field) (*contract.Field, error) {
		if field.IsPK() {
			return nil, nil
		}
		adjusted := field.Optional()
		if err := adjusted.Err(); err != nil {
			return nil, err
		}
		if adjusted.Tag == field.Tag {
			return nil, nil
		}
		return adjusted, nil
	}}, nil
}
func ExampleField_IsPK() {
	key := &contract.Field{Name: "ID", Tag: `sqlx:"ID,primaryKey" required:"true"`}
	ordinary := &contract.Field{Name: "ID", Tag: `json:"id"`}
	fmt.Println(key.IsPK())
	fmt.Println(ordinary.IsPK())
	// Output:
	// true
	// false
}
