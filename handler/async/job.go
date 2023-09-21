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
		ID       string `sqlx:"primaryKey=true,name=ID" json:",omitempty"`
		MatchKey string `json:",omitempty"`
		Status   string `sqlx:"name=Status" json:",omitempty"`
		Metrics  string `json:",omitempty"`
		destination.Table
		destination.Cache
		Request
		Principal
		MainView      string     `json:",omitempty" sqlx:"MainView"`
		Module        string     `json:",omitempty" sqlx:"Module"`
		Labels        string     `json:",omitempty"`
		JobType       string     `json:",omitempty"`
		EventURL      string     `json:",omitempty"`
		Error         *string    `json:",omitempty"`
		CreationTime  time.Time  `json:",omitempty"`
		StartTime     *time.Time `json:",omitempty"`
		EndTime       *time.Time `json:",omitempty"`
		ExpiryTime    *time.Time `json:",omitempty"`
		WaitTimeInMcs int        `json:",omitempty"`
		RunTimeInMcs  int        `json:",omitempty"`
		SQL           []*SQL     `sqlx:"enc=JSON,name=SQLQuery" sqlxAsync:"enc=JSON,name=SQLQuery"`
		Deactivated   bool       `json:",omitempty"`
	}

	Request struct {
		Method string `json:",omitempty"`
		URI    string `json:",omitempty"`
		State  string `json:",omitempty"`
	}

	Principal struct {
		UserEmail *string `json:",omitempty"`
		UserID    *string `json:",omitempty"`
	}

	SQL struct {
		Query string        `json:",omitempty"`
		Args  []interface{} `json:",omitempty"`
	}

	QueryArgs []interface{}

	jobKey string
)

// JobKey defines context job key
var JobKey = jobKey("job")

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
