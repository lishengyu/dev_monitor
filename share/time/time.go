package tf

import "time"

const (
	TimeFormatDefault = "2006-01-02 15:04:05" // 默认时间格式
)

func TimeFormat(t time.Time, format string) string {
	return t.Format(format)
}
