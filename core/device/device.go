package device

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type DeviceInfo struct {
	HostName string // 主机名
	DeviceIP string // IP地址
}

const (
	RouteFile           = "/proc/net/route"
	DefaultRouteFeature = "00000000" // 默认路由标识
)

var DevInfo DeviceInfo

func getDefaultRoute() (string, error) {
	// 读取Linux路由表
	data, err := os.ReadFile(RouteFile)
	if err != nil {
		return "", fmt.Errorf("读取路由表失败: %v", err)
	}

	var defaultInterface string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines[1:] { // 跳过标题行
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[1] == DefaultRouteFeature {
			defaultInterface = fields[0]
			break
		}
	}

	if defaultInterface == "" {
		return "", fmt.Errorf("未找到默认路由接口")
	}
	return defaultInterface, nil
}

func getIPByInterface(ifaceName string) (string, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return "", fmt.Errorf("获取接口失败: %v", err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("获取地址失败: %v", err)
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("接口%s无有效IPv4地址", ifaceName)
}

func GetDeviceIP() (string, error) {
	route, err := getDefaultRoute()
	if err != nil {
		return "", err
	}
	return getIPByInterface(route)
}

func GetHostName() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("获取主机名失败: %v", err)
	}
	return hostName, nil
}

func GetDeviceInfo() *DeviceInfo {
	return &DevInfo
}

func InitDeviceInfo() error {
	var err error
	DevInfo.DeviceIP, err = GetDeviceIP()
	if err != nil {
		return err
	}
	DevInfo.HostName, err = GetHostName()
	if err != nil {
		return err
	}
	return nil
}
