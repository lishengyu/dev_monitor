package tf

import "time"

const (
	LogTimeFormat     = "Jan 2 15:04:05"      // 日志时间格式
	TimeFormatDefault = "2006-01-02 15:04:05" // 默认时间格式
)

func TimeFormat(t time.Time, format string) string {
	return t.Format(format)
}

func TimeFormatIntConvert(t int64, format string) string {
	return time.Unix(t, 0).Format(format)
}
