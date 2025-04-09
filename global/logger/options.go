package logger

import (
	"log/slog"
)

type Options struct {
	Level     slog.Level // 自定义日志等级
	AddSource bool       // 是否记录源代码位置
}

func LogLevelToInt(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func DefaultOptions(level string) *Options {
	return &Options{
		Level:     LogLevelToInt(level),
		AddSource: false,
	}
}
