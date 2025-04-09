package rework

import (
	"dev_monitor/core/rework/draft"
	"dev_monitor/core/sender"
	elog "dev_monitor/global/error"
	"dev_monitor/global/logger"
	"os"
	"regexp"

	"github.com/hpcloud/tail"
)

const (
	SecurityLogFile = "/var/log/security.log"       // 安全日志文件
	SystemLogFile   = "/var/log/messages"           // 系统日志文件
	AuditLogFile    = "/var/log/go-audit/audit.log" // 审计日志文件
)

const (
	SecurityLoginFailedFeature  = "Authentication failure"
	SecurityLoginSuccessFeature = "pam_unix(login:session): session opened for user"
)

var (
	RegexpOnline      *regexp.Regexp
	RegexpOffline     *regexp.Regexp
	RegexpLogin       *regexp.Regexp
	RegexpLogout      *regexp.Regexp
	RegexpLoginFailed *regexp.Regexp
)

func processSecurityLog(text string) error {
	if t, u, ok := draft.MatchUserFeature(RegexpLogin, text); ok {
		oplog, err := draft.NewOperationLogin(t, u)
		if err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
			return err
		}
		sender.ProductData(oplog)
		return nil
	}

	if t, u, ok := draft.MatchUserFeature(RegexpLogout, text); ok {
		oplog, err := draft.NewOperationLogout(t, u)
		if err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
		}
		sender.ProductData(oplog)
	}

	if t, u, ok := draft.MatchUserFeature(RegexpLoginFailed, text); ok {
		oplog, err := draft.NewOperationLoginFail(t, u)
		if err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
		}
		sender.ProductData(oplog)
	}

	return nil
}

// for device online/offline
func processSystemLog(text string) error {
	if t, ok := draft.MatchFeature(RegexpOnline, text); ok {
		oplog, err := draft.NewOperationOffline(t)
		if err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
			return err
		}
		sender.ProductData(oplog)
		return nil
	}

	if t, ok := draft.MatchFeature(RegexpOffline, text); ok {
		oplog, err := draft.NewOperationOnline(t)
		if err != nil {
			logger.Error(elog.StepParse, elog.FailStatus, err)
		}
		sender.ProductData(oplog)
	}

	return nil
}

func processAuditLog(text string) error {
	// 处理审计日志的逻辑
	logger.Debug(elog.StepRawLog, elog.RecordStatus, text)

	audit, err := draft.NewAuditLogByLog(text)
	if err != nil {
		return err
	}

	oplog, err := draft.NewOperationLogByAudit(audit)
	if err != nil {
		return err
	}

	sender.ProductData(oplog)

	return nil
}

func MonitorFile(filename string, f func(line string) error) error {
	var t *tail.Tail
	var err error
	t, err = tail.TailFile(filename, tail.Config{
		Follow:   true,
		ReOpen:   true,
		Logger:   tail.DiscardingLogger,
		Location: &tail.SeekInfo{Whence: os.SEEK_END},
	})
	if err != nil {
		return err
	}

	// 日志处理协程
	go func() {
		logger.Info(elog.StepInit, elog.Routine, elog.Starting, elog.RecordStatus, filename)
		defer logger.Error(elog.StepInit, elog.Routine, elog.Exiting, elog.RecordStatus, filename)
		defer t.Cleanup()

		for {
			select {
			case line := <-t.Lines:
				if line.Err != nil {
					logger.Error(elog.StepParse, elog.FailStatus, line.Err)
					logger.Error(elog.Marshal, elog.FailStatus, err)
					continue
				}
				f(line.Text)
			}
		}
	}()

	return nil
}

// 正则表达式预编译
func InitRegex() {
	RegexpOnline = draft.NewRegexp(draft.OnlineFeature)
	RegexpOffline = draft.NewRegexp(draft.OfflineFeature)
	RegexpLogin = draft.NewRegexp(draft.LoginFeature)
	RegexpLogout = draft.NewRegexp(draft.LogoutFeature)
	RegexpLoginFailed = draft.NewRegexp(draft.LogFailFeature)
}

func MonitorFiles() error {
	var err error
	// 监控audit话单文件
	if err = MonitorFile(AuditLogFile, processAuditLog); err != nil {
		logger.Error(elog.StepInit, elog.FailStatus, err)
		return err
	}
	// 监控安全日志文件
	if err = MonitorFile(SecurityLogFile, processSecurityLog); err != nil {
		logger.Error(elog.StepInit, elog.FailStatus, err)
		return err
	}
	// 监控系统日志文件
	if err = MonitorFile(SystemLogFile, processSystemLog); err != nil {
		logger.Error(elog.StepInit, elog.FailStatus, err)
		return err
	}

	return err
}
