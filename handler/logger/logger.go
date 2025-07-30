package logger

import (
	"context"
	"log/slog"
)

type Logger interface {
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool

	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Infoc(ctx context.Context, msg string, args ...any)
	Debugc(ctx context.Context, msg string, args ...any)
	DebugJSONc(ctx context.Context, msg string, obj any)
	Warnc(ctx context.Context, msg string, args ...any)
	Errorc(ctx context.Context, msg string, args ...any)
	Infos(ctx context.Context, msg string, attrs ...slog.Attr)
	Debugs(ctx context.Context, msg string, attrs ...slog.Attr)
	Warns(ctx context.Context, msg string, attrs ...slog.Attr)
	Errors(ctx context.Context, msg string, attrs ...slog.Attr)
}
