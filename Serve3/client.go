package main

import (
	"bigData/Utils"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	defer conn.Close()
	if err != nil {
		fmt.Println("serve connect error")
		return
	}

	sendDataToServe(conn)
	for {
		receiveDataFromServe(conn)
	}
}

func sendDataToServe(conn net.Conn) {
	if conn == nil {
		fmt.Println("client is empty")
		return
	}
	arr := []string{"a", "b", "c", "d", "e", "f", "abc", "abb"}
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

func receiveDataFromServe(conn net.Conn) {
	if conn == nil {
		fmt.Println("client is empty")
		return
	}
	arr := []string{}
	for {
		buf := make([]byte, 16)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("receive serve data error")
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
			fmt.Println("sort arr", arr)
			arr = nil
		}
	}
}
