package main

import (
	"HuanZhen/config"
	"HuanZhen/core/dnsProxy"
	"HuanZhen/core/portConnCheck"
	"HuanZhen/core/portForwarding"
	"HuanZhen/core/processCheck"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var defaultConfig config.Config

// 读取配置文件
func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("读取配置文件失败")
	}
	err = json.Unmarshal(data, &defaultConfig)
	if err != nil {
		panic("解析配置文件失败")
	}
	//fmt.Println(defaultConfig)

	// 检车配置是否正确，主要检测端口是否冲突，要求每个服务的端口不冲突

	// 启动端口转发服务
	go portForwarding.StartPortForward(defaultConfig)

	// 启动 DNS 代理服务
	go dnsProxy.StartDnsProxy()

	// 启动端口连接检测服务
	go portConnCheck.StartPortConnCheck(defaultConfig.PortConnCheck)

	// 启动进程检测
	go processCheck.StartProcessCheck()

	// 启动数据包检测
	//go pcapCheck.StartCheckPcap()

}

func main() {
	log.Println("@===============⚡幻阵⚡====================@\n ===========虚实之间，防护溯源。=================")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	<-sig
	log.Println("Bye")
}
