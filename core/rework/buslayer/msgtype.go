package buslayer

// 安全等级枚举
type SeverityLevel string

const (
	Emergency     SeverityLevel = "Emergency"
	Alert         SeverityLevel = "Alert"
	Critical      SeverityLevel = "Critical"
	Error         SeverityLevel = "Error"
	Warning       SeverityLevel = "Warning"
	Notice        SeverityLevel = "Notice"
	Informational SeverityLevel = "Informational"
	Debug         SeverityLevel = "Debug"
	UnknownLevel  SeverityLevel = "Unknown"
)

// 操作类型枚举
type OperationType string

const (
	Add          OperationType = "Add"
	Modify       OperationType = "Modify"
	Delete       OperationType = "Delete"
	CopyMove     OperationType = "Copy/Move"
	PowerControl OperationType = "PowerOn/Off/Reboot"
	Login        OperationType = "Login"
	FileTransfer OperationType = "Upload/Download"
	CommandExec  OperationType = "Command"
)

type OperateResult string // 操作结果枚举

const (
	Success OperateResult = "Success"
	Failure OperateResult = "Failure"
)

/*
日志源端IP           SyslogIP           IPV4:D.D.D.D十进制，IPV6：RFC压缩格式的位模式，全大写
源端标识             Facility           16主机层，17网络层，18应用层，19未知其他
设备名称             Hostname           省内唯一标识
日志标识             SyslogSymbol       t:trap(告警)    s:security    u:log
操作员               User               系统登录账号
操作终端IP           TerminalIp         格式同SyslogIP
操作对象             OperationObject    表名、网元名、账号ID、物理端口名
操作类型             OperationType      Add（包括外设备接入）；Modify;Delete（包含外设备弹出）；Copy\Move;Poweron\Poweroff\Reboot;Login(登录堡垒机、远程桌面连接)；Upload\Download文件传输；Command执行命令（如远程关机）
操作开始时间         StartTime          YYYY-MM-DD HH:MM:SS 长格式，24小时制
操作结束时间         EndTime            YYYY-MM-DD HH:MM:SS 长格式，24小时制
安全等级             Severity           Emergency,Alert,Critical,Error,Warning,Notice,Informational,Debug,Unknown等
操作结果             Result             Success;Failure
操作详细信息         OPeration          操作详细参数列表。UTF-8字符串。内容为Josn Object对象格式，使用的所有Key、Key中文名、Value格式（含枚举值中文翻译）需要提供给日志设备厂家
*/
// 主日志结构体
type OperationLog struct {
	SyslogIP        string        `json:"SyslogIP" validate:"ip"`              // 日志源端IP[4](@ref)
	Facility        int           `json:"Facility" validate:"min=16,max=19"`   // 源端标识[4](@ref)
	Hostname        string        `json:"Hostname" validate:"required"`        // 设备名称[4](@ref)
	SyslogSymbol    string        `json:"SyslogSymbol" validate:"oneof=t s u"` // 日志标识[4](@ref)
	User            string        `json:"User"`                                // 操作员[4](@ref)
	TerminalIp      string        `json:"TerminalIp" validate:"ip"`            // 操作终端IP[4](@ref)
	OperationObject string        `json:"OperationObject"`                     // 操作对象[4](@ref)
	OperationType   OperationType `json:"OperationType"`                       // 操作类型[4](@ref)
	StartTime       string        `json:"StartTime"`                           // 操作开始时间[1](@ref)
	EndTime         string        `json:"EndTime"`                             // 操作结束时间[1](@ref)
	Severity        SeverityLevel `json:"Severity"`                            // 安全等级[5](@ref)
	Result          OperateResult `json:"Result" validate:"oneof=Success Failure"`
	Operation       string        `json:"Operation"` // 操作详细信息[4](@ref)
}

// 操作详细信息结构
type OperationDetail struct {
	Parameters  map[string]interface{} `json:"params"`    // 动态参数列表
	KeyMapping  map[string]string      `json:"key_desc"`  // 键值说明
	ValueFormat map[string]string      `json:"value_fmt"` // 值格式规范
}

// 工厂方法创建日志实例
func NewOperationLog() *OperationLog {
	return &OperationLog{}
}
