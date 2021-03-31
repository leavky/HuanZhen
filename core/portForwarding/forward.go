package portForwarding

import (
	"HuanZhen/config"
	"HuanZhen/util"
	"log"
	"strconv"
	"time"
)

// DefaultTimeout is the default timeout period of inactivity for convenience
// sake. It is equivelant to 5 minutes.
const DefaultTimeout = time.Minute * 5

// 同时转发TCP和UDP协议
func StartPortForward(defaultConfig config.Config) {
	//var handle = func(portForwardItem config.PortForward){
	//	// 设置负载均衡
	//	var getforward func() string
	//	var fid = -1
	//	if len(portForwardItem.Forward) > 1 {
	//		getforward = func() string {
	//			fid++
	//			if fid >= len(portForwardItem.Forward) {
	//				fid = 0
	//			}
	//			return portForwardItem.Forward[fid]
	//		}
	//	} else {
	//		getforward = func() string {
	//			return portForwardItem.Forward[0]
	//		}
	//	}
	//
	//	// 处理 Tcp 连接
	//	var doTcpConn = func(conn net.Conn) {
	//		// 处理进来的连接
	//		defer conn.Close()
	//		forw := getforward()
	//		log.Println(forw)
	//		fconn, err := net.Dial("tcp", forw)
	//		if err != nil {
	//			log.Println(err)
	//			return
	//		}
	//		defer fconn.Close()
	//		go io.Copy(conn, fconn)
	//		io.Copy(fconn, conn)
	//	}
	//	// Tcp 监听
	//	Tcplisten, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", portForwardItem.Listen))
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer Tcplisten.Close()
	//	log.Println("[tcp]-Listen on", portForwardItem.Listen)
	//	var handleTcp = func() {
	//		for {
	//			conn, err := Tcplisten.Accept()
	//			if err != nil {
	//				continue
	//			}
	//			log.Println("接收到TCP连接", portForwardItem.Listen)
	//			go doTcpConn(conn)
	//		}
	//	}
	//	go handleTcp()
	//
	//	// TODO：待处理UDP
	//}

	for _, portForwardItem := range defaultConfig.PortForward {
		//log.Println(portForwardItem)
		listenPortList, err := util.ParsePort(portForwardItem.Listen)
		if err != nil {
			log.Println("端口解析错误")
		}
		//log.Println(listenPortList)

		// 监听端口
		// TODO: 待处理TCP协议和负载均衡
		for _, portItem := range listenPortList {
			go StartTcpForward("0.0.0.0:"+strconv.Itoa(portItem), portForwardItem.Forward[0], DefaultTimeout)
			go StartUdpForward("0.0.0.0:"+strconv.Itoa(portItem), portForwardItem.Forward[0], DefaultTimeout)
		}
	}
}
