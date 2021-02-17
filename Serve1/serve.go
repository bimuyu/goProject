package main

import (
	"fmt"
	"net"
)

func main() {
	// 开启一个本店服务器 监听9999端口
	listen, err := net.Listen("tcp", "0.0.0.0:9999")
	defer listen.Close()
	if err != nil {
		fmt.Println("服务器开启失败", err)
		return
	}
	fmt.Println("服务器开启，等待客户端链接")
	// 循环接收客户端请求
	for {
		// 获取客户端链接句柄
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("链接出错", err)
			conn.Close()
			return
		}
		go manageClientConnect(conn)
	}
}

// 处理请求
func manageClientConnect(conn net.Conn) {
	if conn == nil {
		fmt.Println("无效链接")
		return
	}
	// 循环接收请求
	for {
		bytes := make([]byte, 4096)
		// 接收请求数据 并判断接收数据是否正常
		read, err := conn.Read(bytes)
		if read == 0 || err != nil {
			conn.Close()
			return
		}
		fmt.Println("服务端收到：",string(bytes))
		conn.Write(bytes[:read])
	}
}
