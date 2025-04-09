package buslayer

//源端标识
const (
	FacilityHost  = 16 // 主机层
	FacilityNet   = 17 // 网络层
	FacilityApp   = 18 // 应用层
	FacilityOther = 19 // 其他
)

// 日志标识
const (
	Trap     = "t" // 告警
	Security = "s" // 安全
	Log      = "u" // 日志
)

func (op *OperationLog) SetSyslogIP(s string) {
	op.SyslogIP = s
}

func (op *OperationLog) SetTerminalIP(s string) {
	op.TerminalIp = s
}

func (op *OperationLog) SetFacility(s int) {
	op.Facility = s
}

func (op *OperationLog) SetHostName(s string) {
	op.Hostname = s
}

func (op *OperationLog) SetSyslogSymbol(s string) {
	op.SyslogSymbol = s
}

func (op *OperationLog) SetUser(s string) {
	op.User = s
}

func (op *OperationLog) SetOperationObject(s string) {
	op.OperationObject = s
}

func (op *OperationLog) SetOperationType(s OperationType) {
	op.OperationType = s
}

func (op *OperationLog) SetStartTime(s string) {
	op.StartTime = s
}

func (op *OperationLog) SetEndTime(s string) {
	op.EndTime = s
}

func (op *OperationLog) SetSeverity(s SeverityLevel) {
	op.Severity = s
}

func (op *OperationLog) SetResult(s OperateResult) {
	op.Result = s
}

func (op *OperationLog) SetOperation(s string) {
	op.Operation = s
}
