package main

import (
	"bigData/Utils"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	defer conn.Close()
	if err != nil {
		fmt.Println("客户端链接失败")
		return
	}
	time.Sleep(time.Second)
	// 给服务器发送数据
	go func() {
		sendDataToServe(conn)
	}()

	// 接收服务端响应
	manageServeConnect(conn)
}

func manageServeConnect(conn net.Conn) {
	// 循环接收请求
	arr := []int{}
	for {
		bytes := make([]byte, 16)
		// 接收请求数据 并判断接收数据是否正常
		n, err := conn.Read(bytes)
		if err != nil {
			fmt.Println("服务器关闭")
			return
		}
		if n == 16 {
			n1 := Utils.BytesToInt(bytes[:len(bytes)/2])
			n2 := Utils.BytesToInt(bytes[len(bytes)/2:])
			// 定义协议的开头
			if 0 == n1 && 0 == n2 {
				arr = make([]int, 0, 0)
			}
			// 接收数组
			if 1 == n1 {
				arr = append(arr, n2)
			}
			// 定义协议的结尾
			if 0 == n1 && 1 == n2 {
				fmt.Println("排序后的数据", arr)
				arr = nil
			}
		}
	}
}

func sendDataToServe(conn net.Conn) {
	toByte0 := Utils.IntToBytes(0)
	toByte1 := Utils.IntToBytes(1)
	conn.Write(append(toByte0, toByte0...))

	arr := []int{1, 4, 7, 8, 2, 6, 5, 3}
	for i := 0; i < len(arr); i++ {
		tmp := append(toByte1, Utils.IntToBytes(arr[i])...)
		conn.Write(tmp)
	}
	// 结束
	conn.Write(append(toByte0, toByte1...))
}
