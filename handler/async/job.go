package async

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/viant/xdatly/handler/async/destination"
	"reflect"
	"time"
)

var bytesType = reflect.TypeOf([]byte{})

type (
	Job struct {
		JobID   string `sqlx:"primaryKey=true,name=JobID" json:",omitempty"`
		Status  string `sqlx:"name=Status" json:",omitempty"`
		Metrics string `json:",omitempty"`
		Method  string `json:",omitempty"`
		URI     string `json:",omitempty"`
		State   string `json:",omitempty"`
		destination.Table
		destination.Cache
		Request
		Principal
		MainView     string         `json:",omitempty" sqlx:"MainView"`
		Labels       string         `json:",omitempty"`
		JobType      string         `json:",omitempty"`
		Error        *string        `json:",omitempty"`
		CreationTime time.Time      `json:",omitempty"`
		EndTime      *time.Time     `json:",omitempty"`
		TimeTaken    *time.Duration `json:",omitempty"`
		SQL          []*SQL         `sqlx:"enc=JSON,name=SQLQuery" sqlxAsync:"enc=JSON,name=SQLQuery"`
	}

	Request struct {
		Method string `json:",omitempty"`
		URI    string `json:",omitempty"`
		State  string `json:",omitempty"`
	}

	Principal struct {
		UserEmail *string `json:",omitempty"`
		Subject   *string `json:",omitempty"`
	}

	SQL struct {
		Query string        `json:",omitempty"`
		Args  []interface{} `json:",omitempty"`
	}

	QueryArgs []interface{}
)

func (q QueryArgs) Value() (driver.Value, error) {
	marshal, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

func (q QueryArgs) Scan(src any) error {
	asBytes, err := AsString(src)
	if err != nil || len(asBytes) == 0 {
		return err
	}
	return json.Unmarshal(asBytes, &q)
}

func AsString(src any) ([]byte, error) {
	switch actual := src.(type) {
	case string:
		return []byte(actual), nil
	case *string:
		if actual != nil {
			return []byte(*actual), nil
		}

		return nil, nil
	case []byte:
		return actual, nil
	}

	rValue := reflect.ValueOf(src)
	if rValue.Type().ConvertibleTo(bytesType) {
		return rValue.Convert(bytesType).Bytes(), nil
	}

	return nil, fmt.Errorf("unsupported Placeholders database type, expected []byte/string got %T", src)
}
