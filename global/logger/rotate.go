package logger

import (
	"dev_monitor/config"
	"io"
	"os"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

/*
Filename:   "./logs/app.log", // 基础日志文件名（必填）
MaxSize:    500,              // 单文件最大体积（单位：MB）
MaxBackups: 5,                // 保留旧日志文件的最大数量
MaxAge:     30,               // 旧日志保留天数（超过自动删除）
Compress:   true,             // 启用压缩归档（节省磁盘空间）
LocalTime:  true,             // 使用本地时间命名归档文件
*/
func NewRotateLogger(filename string, maxSize, maxBackups, maxAge int, compress, localTime bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  localTime,
	}
}

func NewRotateDefaultLogger(log config.LogConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   log.LogFile,
		MaxSize:    log.MaxSize,
		MaxBackups: log.MaxBackups,
		MaxAge:     log.MaxAge,
		Compress:   log.Compress,
		LocalTime:  log.LocalTime,
	}
}

func NewMultiWriter(log config.LogConfig) io.Writer {
	return io.MultiWriter(NewRotateDefaultLogger(log), os.Stdout)
}
