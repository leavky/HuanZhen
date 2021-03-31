package pcapCheck

import (
	"fmt"
	"github.com/google/gopacket/pcap"
)

func StartCheckPcap() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("获取所有的网路接口失败")
	}
	fmt.Println(devices)
}
