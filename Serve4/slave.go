package main

import (
	"bigData/Utils"
	"fmt"
	"net"
)

func main() {
	serve, err := net.Listen("tcp", "0.0.0.0:9990")
	defer serve.Close()
	if err != nil {
		fmt.Println("server start error", err)
		return
	}
	fmt.Println("start receive data")
	for {
		conn, err := serve.Accept()
		if err != nil {
			fmt.Println("connect error")
			return
		}
		go ManageDataFromReceive(conn)
	}
}

func ManageDataFromReceive(conn net.Conn) {
	if conn == nil {
		fmt.Println("connect error")
		return
	}
	u := new(Utils.QuickSortStruct)
	u.IsFile = false
	u.ReceiveFilePath = ""
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
					u.GetReceiveDataSavePath(true)
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
				fmt.Println("接收数据：",u.Data)
				fmt.Println("start sort data")
				u.IsAsc = true
				if 2 == d3 {
					u.IsAsc = false
				}
				if u.IsFile {
					u.ReceiveFileWrite.Flush()
					u.ReceiveFile.Close()
				}
				u.InitMyFunc()
				u.Sort()
				fmt.Println("finish sort data")
				fmt.Println("返回排序数据：",u.Data)
				fmt.Println("start return data")
				u.ReturnDataFromReceive(conn)
				fmt.Println("finish return data")
			}
		}
	}
}
