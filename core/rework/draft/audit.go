package draft

import (
	"dev_monitor/core/device"
	"dev_monitor/core/rework/buslayer"
	"dev_monitor/core/spec"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"dev_monitor/share/formatstr"
	"dev_monitor/share/parse/kv"
	tf "dev_monitor/share/time"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type AuditLogData struct {
	Sequence  int64             `json:"sequence"`
	Timestamp string            `json:"timestamp"` // 使用字符串类型保留原始格式
	Messages  []AuditLogMessage `json:"messages"`
	UIDMap    map[string]string `json:"uid_map"`
}

type AuditLogMessage struct {
	Type int    `json:"type"`
	Data string `json:"data"` // 原始数据保持字符串格式，后续可按需解析
}

func timeStampToInt64(timestamp string) (int64, error) {
	// 解析时间戳字符串为时间对象
	index := strings.Index(timestamp, ".")
	if index != -1 {
		timestamp = timestamp[:index]
	}

	return strconv.ParseInt(timestamp, 10, 64)
}

func NewAuditLogByLog(line string) (*AuditLogData, error) {
	var audit AuditLogData
	err := json.Unmarshal([]byte(line), &audit)
	if err != nil {
		return &audit, err
	}
	return &audit, nil
}

type AuditMesssagesMap map[int][]map[string]string

func getMessageMap(messages []AuditLogMessage) AuditMesssagesMap {
	messageMap := make(AuditMesssagesMap)
	for _, msg := range messages {
		datas, err := kv.ParseKV(msg.Data)
		if err != nil {
			logger.Warn(elog.StepParse, elog.FailStatus, err)
			continue
		}
		if _, ok := messageMap[msg.Type]; !ok {
			messageMap[msg.Type] = []map[string]string{datas}
		}
		messageMap[msg.Type] = append(messageMap[msg.Type], datas)
	}
	return messageMap
}

func getOperationObject(msgs AuditMesssagesMap) (string, error) {
	if len(msgs) == 0 {
		return "", errors.New("no messages found")
	}

	ms, ok := msgs[spec.AUDIT_SYSCALL]
	if ok {
		for _, m := range ms {
			if _, ok := m["exe"]; ok {
				return m["exe"], nil
			}
		}
	}

	return "", errors.New("no operation object found")
}

func getOperationDetail(msgs AuditMesssagesMap) (string, error) {
	if len(msgs) == 0 {
		return "", errors.New("no messages found")
	}

	ms, ok := msgs[spec.AUDIT_SYSCALL]
	if ok {
		for _, m := range ms {
			if _, ok := m["exe"]; ok {
				data, err := json.Marshal(m)
				if err == nil {
					return string(data), nil
				}
			}
		}
	}

	return "", errors.New("no operation object found")
}

func NewOperationLogByAudit(audit *AuditLogData) (*buslayer.OperationLog, error) {
	oplog := buslayer.NewOperationLog()
	dev := device.GetDeviceInfo()
	oplog.SetSyslogIP(dev.DeviceIP)
	oplog.SetHostName(dev.HostName)
	oplog.SetFacility(buslayer.FacilityHost)
	oplog.SetSyslogSymbol(buslayer.Security)

	user := formatstr.JoinMapWithSep(audit.UIDMap, "/")
	oplog.SetUser(user)
	oplog.SetTerminalIP(dev.DeviceIP)

	msgs := getMessageMap(audit.Messages)
	if len(msgs) == 0 {
		return nil, errors.New("no messages found in audit log")
	}

	opObj, err := getOperationObject(msgs)
	if err != nil {
		return nil, err
	} else {
		oplog.SetOperationObject(opObj)
	}

	second, err := timeStampToInt64(audit.Timestamp)
	if err != nil {
		return nil, err
	} else {
		oplog.SetStartTime(tf.TimeFormatIntConvert(second, tf.TimeFormatDefault))
		oplog.SetEndTime(tf.TimeFormatIntConvert(second, tf.TimeFormatDefault))
	}
	oplog.SetOperationType(buslayer.OperationType(buslayer.Add))
	oplog.SetSeverity(buslayer.SeverityLevel(buslayer.Alert))
	oplog.SetResult(buslayer.Success)

	opDetail, err := getOperationDetail(msgs)
	if err != nil {
		return nil, err
	} else {
		oplog.SetOperation(opDetail)
	}

	return oplog, nil
}
