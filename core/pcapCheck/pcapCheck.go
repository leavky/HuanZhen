package pcapCheck

import (
	"HuanZhen/logger"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"time"
)

var (
	device       string = "\\Device\\NPF_{549147ED-812A-445D-BE58-676F7F53B0D4}"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle

	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	tcpLayer layers.TCP
	udpLayer layers.UDP
	dnsLayer layers.DNS
)

func StartCheckPcap() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		logger.HZLogger.Error("[网络数据包分析模块]-获取网络接口失败")
	}
	for _, device := range devices {
		fmt.Println("\nName", device.Name)
		fmt.Println("\nDescription", device.Description)
		for _, addr := range device.Addresses {
			fmt.Println("-IP:", addr.IP)
			fmt.Println("-Netmask:", addr.Netmask)
			fmt.Println("-Broadaddr:", addr.Broadaddr)

		}
	}

	// open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		logger.HZLogger.Warn("[网络数据包分析模块]-网络接口打开失败-", device)
	}
	defer handle.Close()
	// handle all packet
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
		printPacketInfo(packet)
	}

}

func printPacketInfo(packet gopacket.Packet) {
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet,
		&ethLayer,
		&ipLayer,
		&tcpLayer)

	foundLayerTypes := []gopacket.LayerType{}
	err := parser.DecodeLayers(packet.Data(), &foundLayerTypes)
	if err != nil{
		fmt.Println("decoding layers err", err)
	}

	for _, layerType := range foundLayerTypes{
		if layerType == layers.LayerTypeIPv4{
			fmt.Println("IPV4", ipLayer.SrcIP, "->", ipLayer.DstIP)
		}
		if layerType == layers.LayerTypeTCP{
			fmt.Println("TCP port:", tcpLayer.SrcPort, "->", tcpLayer.DstPort)
			fmt.Println("TCP SYN", tcpLayer.SYN, "| ACK", tcpLayer.Ack)
		}
	}
}
