package main

import (
	"bigData/Utils"
	"fmt"
	"net"
	"sort"
)

func main() {
	serve, err := net.Listen("tcp", "0.0.0.0:9999")
	defer serve.Close()
	if err != nil {
		fmt.Println("serve start error")
		return
	}
	for {
		conn, err := serve.Accept()
		if err != nil {
			fmt.Println("client error")
		}
		go manageClientRequest(conn)
	}
}

func manageClientRequest(conn net.Conn) {
	if conn == nil {
		fmt.Println("client error")
		return
	}
	arr := []string{}
	for {
		buf := make([]byte, 16)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("receive client data error")
			return
		}
		if n != 16 {
			continue
		}
		n1 := Utils.BytesToInt(buf[:len(buf)/2])
		n2 := Utils.BytesToInt(buf[len(buf)/2:])
		// 判断开始
		if 0 == n1 && 0 == n2 {
			arr = make([]string, 0, 0)
		}
		// 开始接收数据
		if 3 == n1 {
			tmp := make([]byte, n2, n2)
			read, _ := conn.Read(tmp)
			if read == n2 {
				arr = append(arr, string(tmp))
			}
		}
		// 判断结束
		if 0 == n1 && 1 == n2 {
			fmt.Println("receive:", arr)
			sort.Strings(arr)
			fmt.Println("sort:", arr)
			// sent to client
			byte0 := Utils.IntToBytes(0)
			byte1 := Utils.IntToBytes(1)
			byte3 := Utils.IntToBytes(3)
			conn.Write(append(byte0, byte0...))
			for i := 0; i < len(arr); i++ {
				// 字符串转字节
				bytes := []byte(arr[i])
				// 字节长度转字节
				lentToBytes := Utils.IntToBytes(len(bytes))
				// 返送数据
				conn.Write(append(append(byte3, lentToBytes...), bytes...))
			}
			conn.Write(append(byte0, byte1...))
		}
	}
}
