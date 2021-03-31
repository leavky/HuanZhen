package portForwarding

import (
	"io"
	"log"
	"net"
	"time"
)

type TcpForwarder struct {
	src     string
	dst     string
	timeout time.Duration // 超时时间
}

func handleConn(conn net.Conn, dst string) {
	client, err := net.Dial("tcp", dst)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		defer conn.Close()
		return
	}
	log.Printf("Forwarding from %v to %v\n", conn.LocalAddr(), client.RemoteAddr())
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func StartTcpForward(src, dst string, timeout time.Duration) {
	var tcpForwarder TcpForwarder
	tcpForwarder.src = src
	tcpForwarder.dst = dst
	tcpForwarder.timeout = timeout

	listener, err := net.Listen("tcp", tcpForwarder.src)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		//TODO: 在此记录恶意的IP地址
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		go handleConn(conn, tcpForwarder.dst)
	}
}
