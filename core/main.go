package main

import (
	"dev_monitor/config"
	"dev_monitor/core/device"
	"dev_monitor/core/dump"
	"dev_monitor/core/rework"
	"dev_monitor/core/sender"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	GoVersion string
	BuildTime string
	Version   string
)

func changeDir() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	if err = os.Chdir(dir); err != nil {
		fmt.Printf("change dir: %s", dir)
		return err
	}

	return nil
}

// 需要额外添加的功能
// 1. 通过命令行查询相关统计信息
// 2. 处理日志信息
// 3. 添加服务信息
// 4. 需要添加命令行查询程序当前运行状态(协程处理状态，routine情况，最新处理时间等)

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Printf("Version    : %s \n", Version)
		fmt.Printf("Build Time : %s \n", BuildTime)
		return
	}

	if err := changeDir(); err != nil {
		fmt.Printf("change dir: %v", err)
		return
	}

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
	logger.Info(elog.StepInit, "Version    :", Version)
	logger.Info(elog.StepInit, "Build Time :", BuildTime)

	// step3 设备基本信息获取
	if err := device.InitDeviceInfo(); err != nil {
		logger.Error(elog.StepDevice, elog.FailStatus, err)
		return
	}

	// step5 正则表达式预编译
	rework.InitRegex()

	// step4 通过syslog将加工的日志内容发送到远程服务器
	if err := sender.StartSyslogSender(config.Cfg.Syslog); err != nil {
		logger.Error(elog.StepInit, elog.FailStatus, err)
		return
	}

	// 加载历史备份记录，确认加载本地/var/log/messages文件，是否触发对应的日志
	if err, ok := dump.LoadDumpFile(); err != nil {
		logger.Error(elog.StepLoadDump, elog.FailStatus, err)
	} else {
		err = rework.LoadLocalSystemLog(ok)
		if err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
		}
	}

	// step6 监控日志文件，对监控到的日志内容进行加工处理
	if err := rework.MonitorFiles(); err != nil {
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
