package sender

import (
	"dev_monitor/core/backup"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"sync"
	"sync/atomic"
	"time"
)

type RetryInfo struct {
	FailCount int
	FailTime  int64
}

type RetryRecords struct {
	maps  sync.Map
	inuse atomic.Int64
}

var RetryRecordMap *RetryRecords

func (m *RetryRecords) LoadOrStore(key string, value *RetryInfo) {
	old, loaded := m.maps.LoadOrStore(key, value)
	if loaded {
		info := old.(*RetryInfo)
		info.FailCount++
		info.FailTime = time.Now().Unix()
	} else {
		value.FailCount = 1
		value.FailTime = time.Now().Unix()
		m.inuse.Add(1)
	}
}

func (m *RetryRecords) Load(key string) (*RetryInfo, bool) {
	if v, ok := m.maps.Load(key); ok {
		return v.(*RetryInfo), true
	}
	return &RetryInfo{}, false
}

func (m *RetryRecords) Delete(key any) {
	m.maps.Delete(key)
	m.inuse.Add(-1)
}

func (m *RetryRecords) InUse() int64 {
	return m.inuse.Load()
}
func StartRetryRoutine(interval int, retry int) {
	RetryRecordMap = new(RetryRecords)
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()

	for range ticker.C {
		RetryRecordMap.maps.Range(func(key, value interface{}) bool {
			data := key.(string)
			info := value.(*RetryInfo)
			if info.FailCount > retry {
				RetryRecordMap.Delete(data)
				backup.FailLogWrite(data)
				return true
			}
			if info.FailTime+int64(interval) <= time.Now().Unix() {
				ProductData(data)
			}
			return true
		})
		logger.Debug(elog.Dynamic, elog.MapCapacity, RetryRecordMap.InUse())
	}
}
