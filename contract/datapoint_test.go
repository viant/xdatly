package contract

import (
	"reflect"
	"testing"
)

func TestFieldDatapoint(t *testing.T) {
	for _, tc := range []struct {
		tag  reflect.StructTag
		want string
	}{
		{`sqlx:"ORDER_ID,primaryKey"`, "ORDER_ID"},
		{`sqlx:"column=ORDER_ID"`, "ORDER_ID"},
		{`sqlx:"OLD,name=ORDER_ID"`, "ORDER_ID"},
		{`sqlx:"ID,column=\"Order ID\""`, "Order ID"},
		{``, "OrderID"}, {`sqlx:""`, "OrderID"},
		{`sqlx:"-"`, "OrderID"}, {`sqlx:"ID,column=-"`, "OrderID"},
		{`sqlx:"ID,column="`, "OrderID"},
		{`sqlx:"ID" on:"ID:ID"`, "OrderID"},
	} {
		field := &Field{Name: "OrderID", Tag: tc.tag}
		if got := field.Datapoint(); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.tag, got, tc.want)
		}
	}
	if (*Field)(nil).Datapoint() != "" || (&Field{}).Datapoint() != "" {
		t.Fatal("nil/empty field has a datapoint")
	}
}

func TestFieldDatapointCommonTagsAllocateNothing(t *testing.T) {
	for _, tag := range []reflect.StructTag{``, `sqlx:"-"`, `sqlx:"ORDER_ID,primaryKey,required=true"`, `sqlx:"column=ORDER_ID"`} {
		field := &Field{Name: "ORDER_ID", Tag: tag}
		if count := testing.AllocsPerRun(1000, func() {
			if field.Datapoint() != "ORDER_ID" {
				t.Fatal("missing datapoint")
			}
		}); count != 0 {
			t.Fatalf("%s: %v allocations", tag, count)
		}
	}
}
