package elog

// 日志状态
const (
	SuccStatus   = "success"
	FailStatus   = "failed"
	RecordStatus = "record"
	WarnStatus   = "warn"
)

// 日志阶段
const (
	StepInit    = "初始化"
	StepConf    = "加载配置文件"
	StepInitLog = "初始化日志模块"
	StepDevice  = "获取设备信息"
	StepRawLog  = "原始日志"
	StepParse   = "解析日志"
	StepFormat  = "格式化日志"
	StepSend    = "发送日志"
)

// 日志动作
const (
	Starting = "staring"
	Exiting  = "exiting"
	Running  = "running"
	Stopping = "stopping"
)

// 功能模块
const (
	Routine = "routine"
	Marshal = "marshal"
	Dynamic = "动态处理信息"
)

// channel统计信息
const (
	ChannelRecv  = "channel recv"
	ChannelSend  = "channel send"
	ChannelUsage = "channel usage"
)
