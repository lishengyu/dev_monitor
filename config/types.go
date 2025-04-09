package config

// 顶层配置结构
type Config struct {
	Base   BaseConfig   `yaml:"Base" json:"base"`
	Log    LogConfig    `yaml:"Log" json:"log"`
	Syslog SyslogConfig `yaml:"Syslog" json:"syslog"`
}

// 基础设备配置
type BaseConfig struct {
	DeviceIP string `yaml:"DeviceIP" json:"deviceIP"` // 设备IP地址，不配置时通过默认路由获取
	HostName string `yaml:"HostName" json:"hostName"` // 主机名，不配置时直接从设备获取
}

// 日志配置
type LogConfig struct {
	LogLevel   string `yaml:"LogLevel" json:"logLevel"`     // debug/info/warn/error
	LogFile    string `yaml:"LogFile" json:"logFile"`       // 日志文件路径
	MaxSize    int    `yaml:"MaxSize" json:"maxSize"`       // 单位MB
	MaxAge     int    `yaml:"MaxAge" json:"maxAge"`         // 保留天数
	MaxBackups int    `yaml:"MaxBackups" json:"maxBackups"` // 保留的旧日志文件的最大数量
	Compress   bool   `yaml:"Compress" json:"compress"`     // 是否启用压缩归档
	LocalTime  bool   `yaml:"LocalTime" json:"localTime"`   // 是否使用本地时间命名归档文件
}

// Syslog服务器配置
type SyslogConfig struct {
	Addr       string `yaml:"Addr" json:"addr"`             // 服务器地址:端口
	Protocol   string `yaml:"Protocol" json:"protocol"`     // tcp/udp
	Tag        string `yaml:"Tag" json:"tag"`               // 日志标签
	ServerCrt  string `yaml:"ServerCrt" json:"serverCrt"`   // TLS证书文件路径
	Priority   int    `yaml:"Priority" json:"priority"`     // 日志级别
	ServerName string `yaml:"ServerName" json:"serverName"` // 证书中的SAN名称，用于TLS验证
}
