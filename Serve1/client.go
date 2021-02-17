package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	dial, err := net.Dial("tcp", "127.0.0.1:9999")
	defer dial.Close()
	if err != nil {
		fmt.Println("客户端链接失败")
		return
	}
	// 接收服务端响应
	go func() {
		manageServeConnect(dial)
	}()

	// 接收客户端输入
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		text := input.Text()
		dial.Write([]byte(text))
	}
}

func manageServeConnect(conn net.Conn) {
	for {
		bytes := make([]byte, 4096)
		read, err := conn.Read(bytes)
		if err != nil {
			fmt.Println("客户端接收响应信息失败")
			return
		}
		fmt.Println("客户端收到：",string(bytes[:read]))
	}
}
