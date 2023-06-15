package logger

type Logger interface {
}

type Service struct {
	logger Logger
}

func NewService(logger Logger) *Service {
	return &Service{logger: logger}
}
