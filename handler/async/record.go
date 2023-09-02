package async

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

var bytesType = reflect.TypeOf([]byte{})

const (
	StateRunning State = "RUNNING"
	StateDone    State = "DONE"
)

const (
	CreateDispositionIfNeeded CreateDisposition = "CREATE_IF_NEEDED"
	CreateDispositionNever    CreateDisposition = "CREATE_NEVER"
)

const (
	WriteDispositionTruncate WriteDisposition = "WRITE_TRUNCATE"
	WriteDispositionEmpty    WriteDisposition = "WRITE_EMPTY"
	WriteDispositionAppend   WriteDisposition = "WRITE_APPEND"
)

type (
	State             string
	CreateDisposition string
	WriteDisposition  string

	JobWithMeta struct {
		Metadata *JobMetadata
		Job      *Job
	}

	JobMetadata struct {
		CacheHit bool
	}

	Job struct {
		JobID   string `sqlx:"primaryKey=true,name=JobID" json:",omitempty"`
		State   State  `sqlx:"name=State" json:",omitempty"`
		Metrics string `json:",omitempty"`
		RecordRequest
		RecordPrincipal
		RecordDestination
		MainView     string         `json:",omitempty" sqlx:"MainView"`
		Labels       string         `json:",omitempty"`
		JobType      string         `json:",omitempty"`
		Error        *string        `json:",omitempty"`
		CreationTime time.Time      `json:",omitempty"`
		EndTime      *time.Time     `json:",omitempty"`
		TimeTaken    *time.Duration `json:",omitempty"`
		SQL          []*SQL         `sqlx:"enc=JSON,name=SQLQuery" sqlxAsync:"enc=JSON,name=SQLQuery"`
	}

	RecordRequest struct {
		RequestRouteURI string `json:",omitempty"`
		RequestURI      string `json:",omitempty"`
		RequestHeader   string `json:",omitempty"`
		RequestMethod   string `json:",omitempty"`
	}

	RecordPrincipal struct {
		PrincipalUserEmail *string `json:",omitempty"`
		PrincipalSubject   *string `json:",omitempty"`
	}

	RecordDestination struct {
		DestinationConnector         string            `json:",omitempty"`
		DestinationDataset           string            `json:",omitempty"`
		DestinationTable             string            `json:",omitempty"`
		DestinationCreateDisposition CreateDisposition `json:",omitempty"`
		DestinationSchema            *string           `json:",omitempty"`
		DestinationTemplate          *string           `json:",omitempty"`
		DestinationWriteDisposition  *WriteDisposition `json:",omitempty"`
		DestinationBucketURL         string            `json:",omitempty"`
		DestinationQueueName         string            `json:",omitempty"`
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
