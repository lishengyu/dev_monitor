package backup

import (
	"dev_monitor/core/backup/rotate"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

//NewRotateWriter

const (
	SuccBackFile = "./back/success/syslog_succ.log" // 成功备份的文件
	FailBackFile = "./back/fail/syslog_fail.log"    // 失败备份的文件
)

type AsyncWriter struct {
	ch          chan []byte
	buffer      []byte
	bufferSize  int
	timeout     time.Duration
	flushLock   sync.Mutex
	rotate      *lumberjack.Logger
	closeSignal chan struct{}
}

var SuccBackWriter *AsyncWriter
var FailBackWriter *AsyncWriter

func NewBackupWriter(filename string) *AsyncWriter {
	return &AsyncWriter{
		ch:         make(chan []byte, 1000),
		bufferSize: 1024 * 1024, // 1MB
		buffer:     make([]byte, 0, 1024*1024),
		timeout:    time.Second * 5,
		rotate:     rotate.NewRotateWriter(filename),
	}
}

func (w *AsyncWriter) Write(data []byte) (int, error) {
	w.ch <- data
	return len(data), nil
}

func (w *AsyncWriter) flush() {
	w.flushLock.Lock()
	defer w.flushLock.Unlock()

	if len(w.buffer) > 0 {
		_, err := w.rotate.Write(w.buffer)
		if err != nil {
			return
		}
		w.buffer = w.buffer[:0] // 清空缓冲区
	}
}

func (w *AsyncWriter) Close() {
	close(w.ch)
	w.flush()
	if err := w.rotate.Close(); err != nil {
		return
	}
}

func (w *AsyncWriter) process() {
	timer := time.NewTimer(w.timeout)
	defer timer.Stop()
	defer logger.Warn(elog.Routine, elog.Exiting, "backup writer")

	for {
		select {
		case data, ok := <-w.ch:
			if !ok {
				w.flush()
				return
			}

			w.buffer = append(w.buffer, data...)
			if len(w.buffer) > w.bufferSize {
				w.flush()
				timer.Reset(w.timeout)
			} else {
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(w.timeout)
			}
		case <-timer.C:
			w.flush()
			timer.Reset(w.timeout)
		case <-w.closeSignal:
			w.flush()
			return
		}
	}
}

func StartSuccBackRoutine(filename string) *AsyncWriter {
	w := NewBackupWriter(filename)
	go w.process()
	return w
}

func StartFailBackRoutine(filename string) *AsyncWriter {
	w := NewBackupWriter(filename)
	go w.process()
	return w
}

func StartBackupLog() {
	SuccBackWriter = StartSuccBackRoutine(SuccBackFile)
	FailBackWriter = StartFailBackRoutine(FailBackFile)
}

func SuccLogWrite(data string) {
	SuccBackWriter.Write([]byte(data + "\n"))
}

func FailLogWrite(data string) {
	FailBackWriter.Write([]byte(data + "\n"))
}
