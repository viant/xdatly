package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
)

type depthKey int

// Setup, context allocation, and handler creation are outside row-loop timings.
// slog's discarded text sink retains enabled formatting cost without terminal IO.
func BenchmarkInvocationLogger(b *testing.B) {
	for _, enabled := range []bool{false, true} {
		level := slog.LevelInfo
		if enabled {
			level = slog.LevelDebug
		}
		log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: level}))
		for _, depth := range []int{0, 4, 16} {
			ctx := WithContext(context.Background(), log)
			for i := 0; i < depth; i++ {
				ctx = context.WithValue(ctx, depthKey(i), i)
			}
			for _, mode := range []string{"direct", "typed-context", "once-per-invocation"} {
				b.Run(fmt.Sprintf("enabled=%t/depth=%d/%s", enabled, depth, mode), func(b *testing.B) {
					b.ReportAllocs()
					b.RunParallel(func(pb *testing.PB) {
						invocationLogger := FromContext(ctx)
						for pb.Next() {
							current := Logger(log)
							switch mode {
							case "typed-context":
								current = FromContext(ctx)
							case "once-per-invocation":
								current = invocationLogger
							}
							current.Debug("begin")
							current.Debug("end")
						}
					})
				})
			}
		}
	}
}

var benchmarkLogger Logger
var benchmarkContext context.Context

func BenchmarkLoggerContextLookup(b *testing.B) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	for _, depth := range []int{0, 4, 16} {
		ctx := WithContext(context.Background(), log)
		for i := 0; i < depth; i++ {
			ctx = context.WithValue(ctx, depthKey(i), i)
		}
		b.Run(fmt.Sprintf("depth=%d", depth), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkLogger = FromContext(ctx)
			}
		})
	}
}
func BenchmarkLoggerContextAttachment(b *testing.B) {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkContext = WithContext(ctx, log)
	}
}
func BenchmarkInvocationLoggerSerial(b *testing.B) {
	for _, enabled := range []bool{false, true} {
		level := slog.LevelInfo
		if enabled {
			level = slog.LevelDebug
		}
		log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: level}))
		for _, depth := range []int{0, 4, 16} {
			ctx := WithContext(context.Background(), log)
			for i := 0; i < depth; i++ {
				ctx = context.WithValue(ctx, depthKey(i), i)
			}
			for _, mode := range []string{"direct", "typed-context", "once-per-invocation"} {
				b.Run(fmt.Sprintf("enabled=%t/depth=%d/%s", enabled, depth, mode), func(b *testing.B) {
					b.ReportAllocs()
					invocationLogger := FromContext(ctx)
					for i := 0; i < b.N; i++ {
						current := Logger(log)
						switch mode {
						case "typed-context":
							current = FromContext(ctx)
						case "once-per-invocation":
							current = invocationLogger
						}
						current.Debug("begin")
						current.Debug("end")
					}
				})
			}
		}
	}
}
