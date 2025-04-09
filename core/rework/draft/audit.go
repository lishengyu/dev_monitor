package draft

import (
	"dev_monitor/core/device"
	"dev_monitor/core/rework/buslayer"
	"dev_monitor/core/spec"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"dev_monitor/share/formatstr"
	"dev_monitor/share/parse/kv"
	"encoding/json"
	"errors"
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

func getOperateionDetail(msgs AuditMesssagesMap) (string, error) {
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

	oplog.SetOperationType(buslayer.OperationType(buslayer.Add))
	oplog.SetStartTime(audit.Timestamp)
	oplog.SetEndTime(audit.Timestamp)
	oplog.SetSeverity(buslayer.SeverityLevel(buslayer.Alert))

	oplog.SetResult(buslayer.Success)

	opDetail, err := getOperateionDetail(msgs)
	if err != nil {
		return nil, err
	} else {
		oplog.SetOperation(opDetail)
	}

	return oplog, nil
}
