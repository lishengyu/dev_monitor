package draft

import (
	"dev_monitor/core/device"
	"dev_monitor/core/rework/buslayer"
	tf "dev_monitor/share/time"
	"encoding/json"
	"regexp"
	"time"
)

/*
Mar 14 18:09:12 localhost sudo: pam_unix(sudo:session): session opened for user root by (uid=0)
Mar 14 18:09:15 localhost sudo: pam_unix(sudo:session): session closed for user root

Apr  2 19:19:45 localhost login: pam_unix(login:auth): authentication failure; logname=LOGIN uid=0 euid=0 tty=tty1 ruser= rhost=  user=root
*/

type OpIndex int

const (
	OpLogin OpIndex = iota
	OpLogout
	OpLoginFail
)

const (
	LoginFeature   = `^(\w{3}\s+\d{1,2}\s\d{2}:\d{2}:\d{2}).*\b(session opened for user)\b(\w+)\b`
	LogoutFeature  = `^(\w{3}\s+\d{1,2}\s\d{2}:\d{2}:\d{2}).*\b(session closed for user)\b(\w+)\b`
	LogFailFeature = `^(\w{3}\s+\d{1,2}\s\d{2}:\d{2}:\d{2}).*\b(authentication failure;)\blogname=(\w+)\b`
)

var (
	OpString = map[OpIndex]string{
		OpLogin:     "login",
		OpLogout:    "logout",
		OpLoginFail: "login failed",
	}

	OpCode = map[OpIndex]string{
		OpLogin:     "success",
		OpLogout:    "success",
		OpLoginFail: "failed",
	}
)

type UserInfo struct {
	Username string `json:"username"` // 用户名
	OpType   string `json:"optype"`   // 操作类型
	OpTime   string `json:"optime"`   // 操作时间
	OpCode   string `json:"opcode"`   // 操作结果
}

// 解析时间，用户名
func MatchUserFeature(re *regexp.Regexp, text string) (time.Time, string, bool) {
	matches := re.FindStringSubmatch(text)
	if len(matches) < 4 {
		return time.Time{}, "", false
	}

	// 解析时间
	parsedTime, err := time.Parse(tf.LogTimeFormat, matches[1])
	if err != nil {
		return time.Time{}, "", false
	}

	return parsedTime, matches[3], true
}

func NewUserInfo(u string, t time.Time, index OpIndex) *UserInfo {
	return &UserInfo{
		Username: u,
		OpType:   OpString[index],
		OpTime:   tf.TimeFormat(t, tf.TimeFormatDefault),
		OpCode:   OpCode[index],
	}
}

func NewOperationUser(t time.Time, user, detail string, result buslayer.OperateResult) (*buslayer.OperationLog, error) {
	oplog := buslayer.NewOperationLog()
	dev := device.GetDeviceInfo()
	oplog.SetSyslogIP(dev.DeviceIP)
	oplog.SetHostName(dev.HostName)
	oplog.SetFacility(buslayer.FacilityHost)
	oplog.SetSyslogSymbol(buslayer.Security)
	oplog.SetUser(user)
	oplog.SetTerminalIP(dev.DeviceIP)
	oplog.SetOperationObject(dev.HostName)
	oplog.SetOperationType(buslayer.OperationType(buslayer.Login))
	oplog.SetStartTime(tf.TimeFormat(t, tf.TimeFormatDefault))
	oplog.SetEndTime(tf.TimeFormat(t, tf.TimeFormatDefault))
	oplog.SetSeverity(buslayer.SeverityLevel(buslayer.Warning))
	oplog.SetResult(result)
	oplog.SetOperation(detail)

	return oplog, nil
}

func NewOperationLogin(t time.Time, user string) (*buslayer.OperationLog, error) {
	detail, err := json.Marshal(NewUserInfo(user, t, OpLogin))
	if err != nil {
		return nil, err
	}
	return NewOperationUser(t, user, string(detail), buslayer.Success)
}

func NewOperationLogout(t time.Time, user string) (*buslayer.OperationLog, error) {
	detail, err := json.Marshal(NewUserInfo(user, t, OpLogout))
	if err != nil {
		return nil, err
	}
	return NewOperationUser(t, user, string(detail), buslayer.Success)
}

func NewOperationLoginFail(t time.Time, user string) (*buslayer.OperationLog, error) {
	detail, err := json.Marshal(NewUserInfo(user, t, OpLoginFail))
	if err != nil {
		return nil, err
	}
	return NewOperationUser(t, user, string(detail), buslayer.Failure)
}
