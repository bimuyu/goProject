package main

import (
	"bigData/Utils"
	"fmt"
	"net"
	"sync"
)

func main() {
	ipList := []string{"127.0.0.1:9990", "127.0.0.1:9991", "127.0.0.1:9992"}
	num := len(ipList)
	// 原始数据路径
	sourcePath := "/Users/magic/web/golang/source/CSDNPWD/csdn_pwd.txt"
	// 结果文件路径
	splitPath := "/Users/magic/web/golang/source/csdn_split/csdn_pwd_"
	// 原始数据文件切割
	path := Utils.SliceFileToSmall(sourcePath, splitPath, num)

	isAsc := false
	isFile := true
	dataType := "string"

	// 创建链接和链接对象
	connList := make([]net.Conn, num, num)
	sortStructList := make([]Utils.QuickSortStruct, num, num)
	for i := 0; i < num; i++ {
		addr, err := net.ResolveTCPAddr("tcp4", ipList[i])
		Utils.ManageError(err)
		conn, err := net.DialTCP("tcp", nil, addr)
		Utils.ManageError(err)
		connList = append(connList, conn)

		// 创建结构体
		u := new(Utils.QuickSortStruct)
		u.IsAsc = isAsc
		u.IsFile = isFile
		u.DataType = dataType
		if isFile {
			u.SendFilePath = path[i]
		}
		sortStructList = append(sortStructList, *u)
	}

	// 文件发送
	for i := 0; i < num; i++ {
		go func() {
			sortStructList[i].SendDataToServe(connList[i])
		}()
	}

	// 接收文件
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			sortStructList[i].ReceiveDataFromServe(connList[i])
			wg.Done()
		}()
	}
	wg.Wait()

	// 文件归并
	// 获取所有文件地址
	pathList := make([]string, num, num)
	for i := 0; i < num; i++ {
		pathList[i] = sortStructList[i].ReceiveFilePath
	}
	// 文件归并
	finalSavePath := "/Users/magic/web/golang/source/csdn_split/"
	one := sortStructList[0].MergeFileListAsOne(pathList, finalSavePath)
	fmt.Println(one)
}
