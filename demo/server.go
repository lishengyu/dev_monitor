// server.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const (
	ServerCrtPath = "./ca/server.crt"
	ServerKeyPath = "./ca/server.key"
	ServerCaPath  = "./ca/ca.crt"
)

func stripPrefix(raw string) string {
	index := strings.Index(raw, "{")
	if index > -1 {
		return raw[index:]
	}
	return raw
}

func formatJsonStr(raw string) (string, error) {
	var obj interface{}

	str := stripPrefix(raw)
	if err := json.Unmarshal([]byte(str), &obj); err != nil {
		return "", err
	}

	formated, err := json.MarshalIndent(obj, "", "    ")
	return string(formated), err
}

func main() {
	// 加载服务端证书
	cert, err := tls.LoadX509KeyPair(ServerCrtPath, ServerKeyPath)
	if err != nil {
		log.Fatal("加载服务端证书失败:", err)
	}

	// 加载客户端CA证书（用于验证客户端）
	clientCAPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(ServerCaPath)
	if err != nil {
		log.Fatal("读取CA证书失败:", err)
	}
	clientCAPool.AppendCertsFromPEM(caCert)

	// TLS配置
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // 强制验证客户端证书
		ClientCAs:    clientCAPool,
		MinVersion:   tls.VersionTLS12,
	}

	// 启动监听
	listener, err := tls.Listen("tcp", ":6514", config)
	if err != nil {
		log.Fatal("监听失败:", err)
	}
	defer listener.Close()

	fmt.Println("Syslog服务端启动，监听端口 6514...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接受连接失败:", err)
			continue
		} else {
			log.Println("接受连接成功")
		}
		defer conn.Close()

		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("读取内容失败:", err)
				break
			}
			str, err := formatJsonStr(string(buf[:n]))
			if err != nil {
				fmt.Println(string(buf[:n]))
			}
			fmt.Println(str)
			fmt.Println("============================")
		}
	}
}
