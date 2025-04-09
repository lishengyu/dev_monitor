package logger

import (
	"dev_monitor/config"
	"fmt"
	"io"
	"log/slog"
	"runtime"
)

const (
	LevelDebug slog.Level = -4
	LevelInfo  slog.Level = 0
	LevelWarn  slog.Level = 4
	LevelError slog.Level = 8
)

var globalLogger *slog.Logger

func ModuleInit(log config.LogConfig) error {
	output := NewMultiWriter(log)
	opts := DefaultOptions(log.LogLevel)
	InitLog(output, opts)
	return nil
}

func InitLog(output io.Writer, opts *Options) {
	baseHandler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		AddSource: opts.AddSource,
		Level:     opts.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 紧凑化时间格式
			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format("2006-01-02 15:04:05 000"))
			}
			return a
		},
	})

	customHandler := NewCustomHandler(baseHandler, opts)
	globalLogger = slog.New(customHandler)
}

// 动态调整日志级别
func SetLevel(level slog.Level) {
	if handler, ok := globalLogger.Handler().(*CustomHandler); ok {
		handler.opts.Level = level
	}
}

func Debug(msg string, args ...any) {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		r := runtime.FuncForPC(pc)
		fn := r.Name()
		_, line := r.FileLine(pc)
		msg = fmt.Sprintf("[%s:%d] %s", fn, line, msg)
	}
	globalLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		r := runtime.FuncForPC(pc)
		fn := r.Name()
		_, line := r.FileLine(pc)
		msg = fmt.Sprintf("[%s:%d] %s", fn, line, msg)
	}
	globalLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		r := runtime.FuncForPC(pc)
		fn := r.Name()
		_, line := r.FileLine(pc)
		msg = fmt.Sprintf("[%s:%d] %s", fn, line, msg)
	}
	globalLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		r := runtime.FuncForPC(pc)
		fn := r.Name()
		_, line := r.FileLine(pc)
		msg = fmt.Sprintf("[%s:%d] %s", fn, line, msg)
	}
	globalLogger.Error(msg, args...)
}
