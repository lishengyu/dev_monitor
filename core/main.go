package main

import (
	"dev_monitor/config"
	"dev_monitor/core/device"
	"dev_monitor/core/rework"
	"dev_monitor/core/sender"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"fmt"
	"time"
)

// 需要额外添加的功能
// 1. 通过命令行查询相关统计信息
// 2. 处理日志信息
// 3. 添加服务信息
// 4. 需要添加命令行查询程序当前运行状态(协程处理状态，routine情况，最新处理时间等)

func main() {
	// step1 配置文件加载
	if err := config.LoadConfig(config.ConfigPath); err != nil {
		fmt.Errorf("加载配置文件失败: %v", err)
		return
	}

	// step2 日志记录模块初始化
	if err := logger.ModuleInit(config.Cfg.Log); err != nil {
		fmt.Errorf("初始化日志模块失败", "failed", err)
	}

	logger.Info(elog.StepInit)

	// step3 设备基本信息获取
	if err := device.InitDeviceInfo(); err != nil {
		logger.Error(elog.StepDevice, elog.FailStatus, err)
		return
	}

	// step4 通过syslog将加工的日志内容发送到远程服务器
	if err := sender.StartSyslogSender(config.Cfg.Syslog); err != nil {
		logger.Error(elog.StepInit, elog.FailStatus, err)
		return
	}

	// step5 正则表达式预编译
	rework.InitRegex()

	// step6 监控日志文件，对监控到的日志内容进行加工处理
	if err := rework.MonitorFiles(); err != nil {
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
