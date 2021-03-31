package portConnCheck

import (
	"HuanZhen/config"
	"HuanZhen/util"
	"log"
	"net"
	"strconv"
)

const bufferSize = 4096

// 用于检测端口连接，原理为监听端口，当检测到端口有连接时，打印远程连接IP
// 还有一种检测方式，就是直接在网卡上检测数据包，检测远程地址
func StartPortConnCheck(conf config.PortConnCheck) {
	log.Println(conf.Ports)
	portList, err := util.ParsePort(conf.Ports)
	if err != nil{
		log.Println("配置文件错误")
	}
	for _, portItem := range portList{
		go portConnCheckTcp("0.0.0.0:"+strconv.Itoa(portItem))
		go portConnCheckUdp("0.0.0.0:"+strconv.Itoa(portItem))
	}
}

func portConnCheckTcp(src string){
	//log.Println("需要监听的端口", src)
	listener, err := net.Listen("tcp", src)
	if err != nil{
		log.Println("[端口连接检测服务]-端口监听失败")
	}
	defer listener.Close()
	for{
			var handleConn = func(conn net.Conn) {
				// 处理进来的连接
				defer conn.Close()
				log.Println("检测到远程连接", conn.LocalAddr(),conn.RemoteAddr())
			}
		conn, err := listener.Accept()
		if err != nil{
			log.Println("[端口连接检测车服务]-端口连接失败")
		}
		go handleConn(conn)
	}
}

func portConnCheckUdp(src string){
	srcAddr, err := net.ResolveUDPAddr("udp", src)
	if err != nil {
		log.Println("[端口连接检测车服务]-端口配置错误")
	}

	listener, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Println("[端口连接检测车服务]-UDP端口监听失败")
	}


	for{
		buf := make([]byte, bufferSize)
		oob := make([]byte, bufferSize)
		_, _, _, addr, err := listener.ReadMsgUDP(buf, oob)
		if err != nil {
			log.Println("forward: failed to read, terminating:", err)
			return
		}

		log.Println("链接到UDP", addr)
	}


	//for {
	//	log.Println("监听到远程UDP连接:", listener.RemoteAddr())
	//}
}
