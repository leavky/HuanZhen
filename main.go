package main

import (
	"HuanZhen/config"
	"HuanZhen/core/dnsProxy"
	"HuanZhen/core/portConnCheck"
	"HuanZhen/core/portForwarding"
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

	go portForwarding.StartPortForward(defaultConfig)
	go dnsProxy.StartDnsProxy()
	go portConnCheck.StartPortConnCheck(defaultConfig.PortConnCheck)
	//go processCheck.StartProcessCheck()

}

func main() {
	log.Println("@===============⚡幻阵⚡====================@")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	<-sig
	log.Println("Bye")
}
