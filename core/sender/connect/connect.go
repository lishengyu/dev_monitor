package connect

import (
	"crypto/tls"
	"crypto/x509"
	"dev_monitor/config"
	"fmt"
	"os"

	"github.com/RackSec/srslog"
)

/*
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
*/

const (
	clientCert = "./ca/client.crt" // PEM格式客户端证书
	clientKey  = "./ca/client.key" // PEM格式客户端私钥
	caCert     = "./ca/ca.crt"     // CA证书
)

// 加载CA证书池
func loadCACertPool(path string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取CA证书失败: %w", err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("解析CA证书失败")
	}

	return pool, nil
}

// 创建TLS配置 (支持双向认证)
func createTLSConfig(san string) (*tls.Config, error) {
	// 加载客户端证书
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("加载客户端证书失败: %w", err)
	}

	// 加载可信CA证书
	caCertPool, err := loadCACertPool(caCert)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert}, // 客户端证书
		RootCAs:      caCertPool,              // 服务端CA验证
		MinVersion:   tls.VersionTLS12,        // 最低TLS1.2
		ServerName:   san,                     // 证书中的SAN名称
	}, nil
}

func NewSysLogTls(cfg config.SyslogConfig) (*srslog.Writer, error) {
	// 1. 初始化TLS配置
	tlsConfig, err := createTLSConfig(cfg.ServerName)
	if err != nil {
		return nil, err
	}

	// 2. 建立Syslog连接
	writer, err := srslog.DialWithTLSConfig(
		cfg.Protocol,                  // 协议类型
		cfg.Addr,                      // 服务器地址
		srslog.Priority(cfg.Priority), // 设备类型
		cfg.Tag,                       // 应用标签
		tlsConfig,                     // TLS配置
	)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	// 3. 配置日志格式
	writer.SetFormatter(srslog.RFC5424Formatter) // 使用RFC5424标准
	/*
		4. 发送测试日志
		for {
			err = writer.Info("this is a test")
			if err != nil {
				fmt.Println("write failed: ", err)
			}
			time.Sleep(5 * time.Second)
		}
	*/

	return writer, nil
}

func WriterStop(w *srslog.Writer) {
	// 关闭连接
	if w != nil {
		w.Close()
	}
}
