package rotate

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewRotateWriter(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // MB
		MaxBackups: 30,
		MaxAge:     30, // days
		Compress:   true,
		LocalTime:  true,
	}
}
