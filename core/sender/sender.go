package sender

import (
	"dev_monitor/config"
	"dev_monitor/core/sender/connect"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/RackSec/srslog"
)

type MonitorChannel struct {
	ch        chan interface{}
	SendCount atomic.Int64
	RecvCount atomic.Int64
	Usage     atomic.Int64
}

var MonitorChannelInstance *MonitorChannel

func NewMonitorChannel() *MonitorChannel {
	return &MonitorChannel{
		ch: make(chan interface{}, 1000),
	}
}

func (mc *MonitorChannel) Send(data interface{}) {
	mc.ch <- data
	mc.SendCount.Add(1)
}

func (mc *MonitorChannel) Recv() interface{} {
	data := <-mc.ch
	mc.RecvCount.Add(1)
	mc.UsageStore(len(mc.ch))
	return data
}

func (mc *MonitorChannel) UsageStore(use int) {
	usage := int64(use * 100 / cap(mc.ch))
	if usage > 80 {
		logger.Warn(elog.Dynamic,
			elog.ChannelSend, mc.RecvCount.Load(),
			elog.ChannelRecv, mc.RecvCount.Load(),
			elog.ChannelUsage, usage)
	}
	mc.Usage.Store(usage)
}

func ProductData(data interface{}) {
	logger.Debug(elog.StepFormat, elog.RecordStatus, data)
	MonitorChannelInstance.Send(data)
}

func ConsumData(sr *srslog.Writer) {
	for data := range MonitorChannelInstance.ch {
		if err := process(sr, data); err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
		}
	}
}

func ConvertToString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	case error:
		return x.Error()
	case int:
		return strconv.Itoa(x)
	case bool:
		return strconv.FormatBool(x)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func SysLogMarshal(a interface{}) (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

func process(sr *srslog.Writer, data interface{}) error {
	ctx, err := SysLogMarshal(data)
	if err != nil {
		logger.Error(elog.Marshal, elog.FailStatus, err)
		return err
	}

	// 发送数据到Syslog服务器
	err = sr.Info(ctx)
	if err != nil {
		logger.Error(elog.StepSend, elog.FailStatus, err)
		//todo 重试机制 + 错误处理 + 缓存机制
		return err
	}
	logger.Info(elog.StepSend, elog.SuccStatus, ctx)
	return nil
}

func StartConsumer(sr *srslog.Writer) {
	// 处理数据流中的数据
	go func() {
		logger.Info(elog.StepInit, elog.Routine, elog.Starting)
		defer logger.Error(elog.StepInit, elog.Routine, elog.Exiting)
		defer sr.Close()
		ConsumData(sr)
	}()
}

func StartSyslogSender(cfg config.SyslogConfig) error {
	MonitorChannelInstance = NewMonitorChannel()
	sr, err := connect.NewSysLogTls(cfg)
	if err != nil {
		return err
	}

	StartConsumer(sr)
	return nil
}
