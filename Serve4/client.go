package main

import (
	"bigData/Utils"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9990")
	defer conn.Close()
	if err != nil {
		fmt.Println("serve connect error")
		return
	}

	u := new(Utils.QuickSortStruct)

	sendDataToServe(conn, *u)
	for {
		receiveDataFromServe(conn, *u)
	}
}

func sendDataToServe(conn net.Conn, u Utils.QuickSortStruct) {
	if conn == nil {
		fmt.Println("client is empty")
		return
	}

	arr := []string{"a", "b", "c", "d", "e", "f", "abc", "abb"}
	u.IsAsc = true
	u.IsFile = true
	u.DataType = "string"
	//arr := []int{7, 9, 2, 8, 3, 3, 3, 9, 9, 11, 17, 16, 13}
	// sent to client
	byte0 := Utils.IntToBytes(0)
	byte1 := Utils.IntToBytes(1)

	var engine []byte
	if u.IsFile {
		engine = byte1
	} else {
		engine = byte0
	}
	var sortByte []byte
	if u.IsAsc {
		sortByte = byte1
	} else {
		sortByte = Utils.IntToBytes(2)
	}
	start := append(append(engine, byte0...), byte0...)
	conn.Write(start)
	for i := 0; i < len(arr); i++ {
		toByte := u.DataFormatToByte(arr[i])
		conn.Write(append(engine, toByte...))
	}
	end := append(append(engine, byte0...), sortByte...)
	conn.Write(end)
}

func receiveDataFromServe(conn net.Conn, u Utils.QuickSortStruct) {
	if conn == nil {
		fmt.Println("connect error")
		return
	}
	for {
		byt := make([]byte, 24)
		n, err := conn.Read(byt)
		if err != nil {
			fmt.Println("connect close")
			return
		}
		if n == 24 {
			// 第一个字节表示内存还是文件模式 100 000
			d1 := Utils.BytesToInt(byt[:8])
			d2 := Utils.BytesToInt(byt[8:16])
			d3 := Utils.BytesToInt(byt[16:])
			// 存储文件排序
			if d1 == 1 {
				u.IsFile = true
				// 获取文件存储位置
				if u.ReceiveFilePath == "" {
					u.GetReceiveDataSavePath(false)
				}
			}
			// 表示开始数据传输
			if d2 == 0 && d3 == 0 {
				fmt.Println("start receive data")
				u.Data = make([]interface{}, 0, 0)
			}
			// 接收对应的数据类型
			if d2 == 1 { // 整数
				u.DataType = "int"
				u.PutEveryReceiveToData(d3)
			} else if d2 == 2 { // 浮点数
				u.DataType = "float64"
				u.PutEveryReceiveToData(Utils.BytesToFloat64(byt[16:]))
			} else if d2 == 3 { // 字符串
				strByt := make([]byte, d3, d3)
				read, _ := conn.Read(strByt)
				if read == d3 {
					u.DataType = "string"
					u.PutEveryReceiveToData(string(strByt))
				}
			} else if d2 == 4 { // 结构体
				jsonByt := make([]byte, d3, d3)
				read, _ := conn.Read(jsonByt)
				if read == d3 {
					u.DataType = "struct"
					u.PutEveryReceiveToData(jsonByt)
				}
			}
			// 数据接收结束
			if d2 == 0 && (d3 == 1 || d3 == 2) {
				fmt.Println("finish receive data")
				if u.IsFile {
					u.ReceiveFileWrite.Flush()
					fmt.Println(u.ReceiveFilePath)
				} else {
					fmt.Println(u.Data)
				}
			}
		}
	}
}
