package main

import (
	"bigData/Utils"
	"fmt"
	"net"
	"sort"
)

// 定一个简单的整数数组接收协议
// 接收数组排序后返回
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
			fmt.Println("连接出错", err)
		}
		go manageClientConnect(conn)
	}
}

// 处理请求
func manageClientConnect(conn net.Conn) {
	if conn == nil {
		fmt.Println("无效连接")
		return
	}
	// 循环接收请求
	arr := []int{}
	for {
		buf := make([]byte, 16)
		// 接收请求数据 并判断接收数据是否正常
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("客户端关闭")
			return
		}
		if n == 16 {
			n1 := Utils.BytesToInt(buf[:len(buf)/2])
			n2 := Utils.BytesToInt(buf[len(buf)/2:])
			// 定义协议的开头
			if n1 == 0 && n2 == 0 {
				arr = make([]int, 0, 0)
			}
			// 接收数组
			if n1 == 1 {
				arr = append(arr, n2)
			}
			// 定义协议的结尾
			if n1 == 0 && n2 == 1 {
				fmt.Println("收到数组：", arr)
				sort.Ints(arr)
				fmt.Println("收到数组排序完成：", arr)

				myArr := arr

				toByte0 := Utils.IntToBytes(0)
				toByte1 := Utils.IntToBytes(1)
				conn.Write(append(toByte0, toByte0...))

				for i := 0; i < len(myArr); i++ {
					conn.Write(append(toByte1, Utils.IntToBytes(myArr[i])...))
				}

				conn.Write(append(toByte0, toByte1...))
			}
		}
	}
}
