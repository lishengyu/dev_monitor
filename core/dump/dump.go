package dump

import (
	tf "dev_monitor/share/time"
	"encoding/json"
	"os"
	"time"
)

type DumpData struct {
	LogTime string `json:"log_time"`
	Module  string `json:"module"`
}

const (
	DumpFilePath = "./config/dump.json"
)

var DumpRecord DumpData

func DumpFile(t time.Time) error {
	dump := DumpData{
		LogTime: tf.TimeFormat(t, tf.TimeFormatDefault),
		Module:  "system",
	}

	data, err := json.Marshal(dump)
	if err != nil {
		return err
	}

	err = os.WriteFile(DumpFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadDumpFile() (error, bool) {
	if _, err := os.Stat(DumpFilePath); os.IsNotExist(err) {
		return nil, false
	}

	data, err := os.ReadFile(DumpFilePath)
	if err != nil {
		return err, true
	}

	err = json.Unmarshal(data, &DumpRecord)
	if err != nil {
		return err, true
	}

	return nil, true
}

func GetDumpRecordTime() string {
	return DumpRecord.LogTime
}

func AfterDumpTime(t time.Time) bool {
	if t.IsZero() {
		return false
	}

	tr, err := time.Parse(tf.LogTimeFormat, DumpRecord.LogTime)
	if err != nil {
		return false
	}

	if t.After(tr) {
		return true
	}

	return false
}
