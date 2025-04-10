package draft

import (
	"dev_monitor/core/device"
	"dev_monitor/core/rework/buslayer"
	tf "dev_monitor/share/time"
	"regexp"
	"time"
)

/*
Apr  2 20:13:34 localhost systemd-logind: Power key pressed.
Apr  2 20:13:34 localhost systemd-logind: Powering Off...
Apr  2 20:13:34 localhost systemd-logind: System is powering down.
*/

/*
Apr  3 09:55:25 localhost kernel: Initializing cgroup subsys cpuset
Apr  3 09:55:25 localhost kernel: Initializing cgroup subsys cpu
Apr  3 09:55:25 localhost kernel: Initializing cgroup subsys cpuacct
Apr  3 09:55:25 localhost kernel: Linux version 3.10.0-514.el7.x86_64 (builder@kbuilder.dev.centos.org) (gcc version 4.8.5 20150623 (Red Hat 4.8.5-11) (GCC) ) #1 SMP Tue Nov 22 16:42:41 UTC 2016
Apr  3 09:55:25 localhost kernel: Command line: BOOT_IMAGE=/vmlinuz-3.10.0-514.el7.x86_64 root=/dev/mapper/cl-root ro isolcpus=1-110 transparent_hugepage=never default_hugepagesz=1G hugepagesz=1G crashkernel=auto rd.lvm.lv=cl/root rd.lvm.lv=cl/swap nomodeset rhgb quiet
*/

const (
	OnlineFeature  = `^(\w{3}\s+\d{1,2}\s\d{2}:\d{2}:\d{2}).*\b(systemd-logind: System is powering down)\b`
	OfflineFeature = `^(\w{3}\s+\d{1,2}\s\d{2}:\d{2}:\d{2}).*\b(kernel: Command line: BOOT_IMAGE=)\b`
	OnlineDetail   = `{"device op": "online"}`
	OfflineDetail  = `{"device op": "offline"}`
)

func NewRegexp(restr string) *regexp.Regexp {
	return regexp.MustCompile(restr)
}

func MatchFeature(re *regexp.Regexp, text string) (time.Time, bool) {
	matches := re.FindStringSubmatch(text)
	if len(matches) < 3 {
		return time.Time{}, false
	}

	// 解析时间
	parsedTime, err := time.Parse(tf.LogTimeFormat, matches[1])
	if err != nil {
		return time.Time{}, false
	}

	return parsedTime, true
}

func NewOperationPower(t time.Time, detail string) (*buslayer.OperationLog, error) {
	oplog := buslayer.NewOperationLog()
	dev := device.GetDeviceInfo()
	oplog.SetSyslogIP(dev.DeviceIP)
	oplog.SetHostName(dev.HostName)
	oplog.SetFacility(buslayer.FacilityHost)
	oplog.SetSyslogSymbol(buslayer.Trap)
	oplog.SetUser(UserDefault)
	oplog.SetTerminalIP(dev.DeviceIP)
	oplog.SetOperationObject(dev.HostName)
	oplog.SetOperationType(buslayer.OperationType(buslayer.PowerControl))
	oplog.SetStartTime(tf.TimeFormat(t, tf.TimeFormatDefault))
	oplog.SetEndTime(tf.TimeFormat(t, tf.TimeFormatDefault))
	oplog.SetSeverity(buslayer.SeverityLevel(buslayer.Emergency))
	oplog.SetResult(buslayer.Success)
	oplog.SetOperation(detail)

	return oplog, nil
}

func NewOperationOffline(t time.Time) (*buslayer.OperationLog, error) {
	return NewOperationPower(t, OnlineDetail)
}

func NewOperationOnline(t time.Time) (*buslayer.OperationLog, error) {
	return NewOperationPower(t, OfflineDetail)
}
