package main

import (
	"bigData/Utils"
	"fmt"
)

func main() {
	// 原始数据路径
	sourcePath := "/Users/magic/web/golang/source/CSDNPWD/csdn_pwd.txt"
	// 结果文件路径
	savePath := "/Users/magic/web/golang/source/csdn_split/csdn_pwd_"
	// 原始数据文件切割
	path := Utils.SliceFileToSmall(sourcePath, savePath, 9)
	fmt.Println(path)
	// 切割文件发送

	// 接收文件

	// 文件归并
}
