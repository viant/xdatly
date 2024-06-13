package http

type ErrorStatusCoder interface {
	ErrorStatusCode() int
}

type ErrorMessager interface {
	ErrorMessage() string
}

type ErrorObjecter interface {
	ErrorObject() interface{}
}

type Error struct {
	View       string      `json:",omitempty" default:"nullable=true,required=false,allowEmpty=true"`
	Parameter  string      `json:",omitempty" default:"nullable=true,required=false,allowEmpty=true"`
	StatusCode int         `json:",omitempty" default:"nullable=true,required=false,allowEmpty=true"`
	Err        error       `json:"-"`
	Message    string      `json:",omitempty" default:"nullable=true,required=false,allowEmpty=true"`
	Object     interface{} `json:",omitempty" default:"nullable=true,required=false,allowEmpty=true"`
}

func (e *Error) ErrorStatusCode() int {
	return e.StatusCode
}

func (e *Error) ErrorMessage() string {
	return e.Message
}

func (e *Error) ErrorObject() interface{} {
	return e.Object
}
func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message string, object interface{}) *Error {
	return &Error{
		StatusCode: code,
		Message:    message,
		Object:     object,
	}
}
