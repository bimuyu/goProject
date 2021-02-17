package Utils

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 服务器接收文件路径
const SERVICE_FILE_PATH = "/Users/magic/web/golang/source/service/"

// 客户端接收文件路径
const CLIENT_FILE_PATH = "/Users/magic/web/golang/source/client/"

type QuickSortStruct struct {
	Data []interface{}
	// 排序顺序
	IsAsc bool
	// 数据类型
	DataType string
	// 内存还是文件
	IsFile bool
	// 接收文件数据存储位置
	ReceiveFilePath string
	// 接收文件数据存储位置
	ReceiveFileWrite *bufio.Writer
	ReceiveFile      *os.File
	// 要返回数据存储位置
	ReturnFilePath string
	// 要发送的数据路径
	SendFilePath string
	// 函数指针
	MyFunc func(i, j interface{}) bool
}

// 排序
func (this *QuickSortStruct) Sort() {
	// 从文件读入数据
	if this.IsFile && this.ReceiveFilePath != "" {
		this.LoadDataFromFile()
	}
	// 排序
	if len(this.Data) <= 100 {
		this.BinSearchSort()
	} else {
		this.QuickSort(0, len(this.Data)-1)
	}
	fmt.Println(this.Data)
	// 写入数据
	if this.IsFile && this.ReceiveFilePath != "" {
		this.SaveSortDataToFile()
	}
}

// 二分插入排序
func (this *QuickSortStruct) BinSearchSort() {
	if len(this.Data) <= 1 {
		return
	} else {
		for i := 1; i < len(this.Data); i++ {
			index := this.FindInsertIndex(0, i-1, i)
			if i != index {
				for j := i; j > index; j-- {
					this.Data[j], this.Data[j-1] = this.Data[j-1], this.Data[j]
				}
			}
		}
	}
}

// 交换位置
func (this *QuickSortStruct) Swap(i, j int) {
	this.Data[i], this.Data[j] = this.Data[j], this.Data[i]
}

// 获取要插入的位置
func (this *QuickSortStruct) FindInsertIndex(start, end, current int) int {
	if start >= end {
		if this.MyFunc(this.Data[start], this.Data[current]) {
			return current
		} else {
			return start
		}
	}
	mid := (start + end) / 2
	if this.MyFunc(this.Data[mid], this.Data[current]) {
		return this.FindInsertIndex(mid+1, end, current)
	} else {
		return this.FindInsertIndex(start, mid, current)
	}
}

// 快速排序
func (this *QuickSortStruct) QuickSort(start, end int) {
	if end-start <= 100 {
		this.SortPartOfData(start, end)
	} else {
		this.Swap(start, rand.Int()%(end-start+1)+start)
		tmp := this.Data[start]
		lt := start
		gt := end + 1
		i := start + 1
		for i < gt {
			if this.MyFunc(this.Data[i], tmp) {
				lt++
				this.Swap(i, lt)
				i++
			} else if this.MyFunc(tmp, this.Data[i]) {
				gt--
				this.Swap(i, gt)
			} else {
				i++
			}
		}
		this.Swap(start, lt)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			this.QuickSort(start, lt-1)
			wg.Done()
		}()
		go func() {
			this.QuickSort(gt, end)
			wg.Done()
		}()
		wg.Wait()
	}
}

// 截取数组的部分长度进行排序
func (this *QuickSortStruct) SortPartOfData(start, end int) {
	if end-start <= 1 {
		if !this.MyFunc(this.Data[start], this.Data[end]) {
			this.Swap(start, end)
		}
	} else {
		for i := start + 1; i <= end; i++ {
			index := this.FindInsertIndex(start, i-1, i)
			if i != index {
				for j := i; j > index; j-- {
					this.Data[j], this.Data[j-1] = this.Data[j-1], this.Data[j]
				}
			}
		}
	}
}

// 接收的数据依次放入结构体data中
func (this *QuickSortStruct) PutEveryReceiveToData(data interface{}) {
	if data == nil {
		fmt.Println("empty data")
		return
	}
	if "int" == this.DataType {
		tmp := data.(int)
		if !this.IsFile {
			this.Data = append(this.Data, tmp)
		} else {
			this.ReceiveFileWrite.WriteString(strconv.Itoa(tmp) + "\n")
		}
	} else if "float64" == this.DataType {
		tmp := data.(float64)
		if !this.IsFile {
			this.Data = append(this.Data, tmp)
		} else {
			this.ReceiveFileWrite.WriteString(strconv.FormatFloat(tmp, 'f', 2, 64) + "\n")
		}
	} else if "string" == this.DataType {
		tmp := data.(string)
		if !this.IsFile {
			this.Data = append(this.Data, tmp)
		} else {
			this.ReceiveFileWrite.WriteString(tmp + "\n")
		}
	} else if "struct" == this.DataType {
		tmp := new(SortDemoStruct)
		err := json.Unmarshal(data.([]byte), &tmp)
		if err != nil {
			fmt.Println("json unmarshal err", err)
			return
		}
		if !this.IsFile {
			this.Data = append(this.Data, tmp)
		} else {
			this.ReceiveFileWrite.WriteString(tmp.Info + " # " + strconv.Itoa(tmp.Times) + "\n")
		}
	} else {
		fmt.Println("unknow data type")
		return
	}
}

// 返回发送来的数据
func (this *QuickSortStruct) ReturnDataFromReceive(conn net.Conn) {
	var sortByte []byte
	if this.IsAsc {
		sortByte = IntToBytes(1)
	} else {
		sortByte = IntToBytes(2)
	}
	byte0 := IntToBytes(0)
	if !this.IsFile { // 内存模式
		start := append(append(byte0, byte0...), byte0...)
		conn.Write(start)
		for i := 0; i < len(this.Data); i++ {
			byt := this.DataFormatToByte(this.Data[i])
			if byt == nil {
				continue
			}
			tmp := append(byte0, byt...)
			conn.Write(tmp)
		}
		end := append(append(byte0, byte0...), sortByte...)
		conn.Write(end)
	} else { // 文件模式
		byte1 := IntToBytes(1)
		start := append(append(byte1, byte0...), byte0...)
		conn.Write(start)
		open, err := os.Open(this.ReturnFilePath)
		defer open.Close()
		if err != nil {
			fmt.Println("ReturnFilePath open err")
			return
		}
		reader := bufio.NewReader(open)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			toByte := this.FileDataFormatToByte(string(line))
			tmp := append(byte1, toByte...)
			conn.Write(tmp)
		}
		end := append(append(byte1, byte0...), sortByte...)
		conn.Write(end)
		open.Close()
	}
}

// 创建要保存数据路径
func (this *QuickSortStruct) GetReceiveDataSavePath(isService bool) {
	for {
		nano := time.Now().UnixNano() + int64(rand.Int())
		path := ""
		if isService {
			path = SERVICE_FILE_PATH + strconv.Itoa(int(nano)) + ".txt"
		} else {
			path = CLIENT_FILE_PATH + strconv.Itoa(int(nano)) + ".txt"
		}
		_, err := os.Stat(path)
		if err != nil {
			create, err := os.Create(path)
			fmt.Println(err)
			if err != nil {
				fmt.Println("file create error", err)
				os.Exit(404)
				return
			}
			writer := bufio.NewWriter(create)
			this.ReceiveFilePath = path
			this.ReceiveFileWrite = writer
			this.ReceiveFile = create
			return
		}
	}
}

// 从文件加载数据
func (this *QuickSortStruct) LoadDataFromFile() {
	open, err := os.Open(this.ReceiveFilePath)
	if err != nil {
		fmt.Println("file open err")
		return
	}
	reader := bufio.NewReader(open)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		if "int" == this.DataType {
			atoi, _ := strconv.Atoi(tmp)
			this.Data = append(this.Data, atoi)
		} else if "float64" == this.DataType {
			float, _ := strconv.ParseFloat(tmp, 64)
			this.Data = append(this.Data, float)
		} else if "string" == this.DataType {
			this.Data = append(this.Data, tmp)
		} else {
			split := strings.Split(tmp, " # ")
			s := new(SortDemoStruct)
			s.Info = split[0]
			atoi, _ := strconv.Atoi(split[1])
			s.Times = atoi
			this.Data = append(this.Data, s)
		}
	}
	open.Close()
}

// 排序数据写入文件
func (this *QuickSortStruct) SaveSortDataToFile() {
	recvFile := this.ReceiveFilePath
	split := strings.Split(recvFile, "/")
	fileName := split[len(split)-1]
	saveName := strings.Replace(fileName, ".txt", "_sort.txt", -1)
	savePath := SERVICE_FILE_PATH + saveName
	open, err := os.Create(savePath)
	defer open.Close()
	if err != nil {
		fmt.Println("save path create err")
		return
	}
	this.ReturnFilePath = savePath
	writer := bufio.NewWriter(open)
	for i := 0; i < len(this.Data); i++ {
		if "int" == this.DataType {
			fmt.Fprintln(writer, strconv.Itoa(this.Data[i].(int)))
		} else if "float64" == this.DataType {
			float := strconv.FormatFloat(this.Data[i].(float64), 'f', 2, 64)
			fmt.Fprintln(writer, float)
		} else if "string" == this.DataType {
			fmt.Fprintln(writer, this.Data[i].(string))
		} else {
			tmp := this.Data[i].(SortDemoStruct)
			fmt.Fprintln(writer, tmp.Info+" # "+strconv.Itoa(tmp.Times))
		}
	}
	writer.Flush()
	this.Data = nil
	runtime.GC()
	debug.FreeOSMemory()
}

// 数据转为字节
func (this *QuickSortStruct) DataFormatToByte(data interface{}) []byte {
	var byt []byte
	if "int" == this.DataType {
		byt = append(append(byt, IntToBytes(1)...), IntToBytes(data.(int))...)
	} else if "float64" == this.DataType {
		byt = append(append(byt, IntToBytes(2)...), Float64ToBytes(data.(float64))...)
	} else if "string" == this.DataType {
		tmp := []byte(data.(string))
		lengthByt := IntToBytes(len(tmp))
		byt = append(append(byt, IntToBytes(3)...), lengthByt...)
		byt = append(byt, tmp...)
	} else if "struct" == this.DataType {
		tmp := data.(SortDemoStruct)
		marshal, err := json.Marshal(tmp)
		if err != nil {
			fmt.Println("json format err", err)
			return nil
		}
		lengthByt := IntToBytes(len(marshal))
		byt = append(append(byt, IntToBytes(4)...), lengthByt...)
		byt = append(byt, marshal...)
	} else {
		return nil
	}
	return byt
}

// 文件数据转为字节
func (this *QuickSortStruct) FileDataFormatToByte(data string) []byte {
	var byt []byte
	var dataTypeByte []byte
	var dataByte []byte
	if "int" == this.DataType {
		atoi, _ := strconv.Atoi(data)
		dataTypeByte = IntToBytes(1)
		dataByte = IntToBytes(atoi)
	} else if "float64" == this.DataType {
		dataTypeByte = IntToBytes(2)
		float, _ := strconv.ParseFloat(data, 64)
		dataByte = Float64ToBytes(float)
	} else if "string" == this.DataType {
		dataTypeByte = IntToBytes(3)
		tmp := []byte(data)
		lengthByt := IntToBytes(len(tmp))
		dataByte = append(lengthByt, tmp...)
	} else if "struct" == this.DataType {
		dataTypeByte = IntToBytes(4)
		split := strings.Split(data, " # ")
		s := new(SortDemoStruct)
		s.Info = split[0]
		atoi, _ := strconv.Atoi(split[1])
		s.Times = atoi
		marshal, err := json.Marshal(s)
		if err != nil {
			fmt.Println("json marshal error")
			return nil
		}
		lengthByt := IntToBytes(len(marshal))
		dataByte = append(lengthByt, marshal...)
	} else {
		return nil
	}
	return append(append(byt, dataTypeByte...), dataByte...)
}

// 初始化myFunc
func (this *QuickSortStruct) InitMyFunc() {
	if "int" == this.DataType {
		this.MyFunc = func(i, j interface{}) bool {
			if this.IsAsc {
				return i.(int) < j.(int)
			} else {
				return i.(int) > j.(int)
			}
		}
	} else if "float64" == this.DataType {
		this.MyFunc = func(i, j interface{}) bool {
			if this.IsAsc {
				return i.(float64) < j.(float64)
			} else {
				return i.(float64) > j.(float64)
			}
		}
	} else if "string" == this.DataType {
		this.MyFunc = func(i, j interface{}) bool {
			if this.IsAsc {
				return i.(string) < j.(string)
			} else {
				return i.(string) > j.(string)
			}
		}
	} else if "struct" == this.DataType {
		this.MyFunc = func(i, j interface{}) bool {
			if this.IsAsc {
				return i.(SortDemoStruct).Times < j.(SortDemoStruct).Times
			} else {
				return i.(SortDemoStruct).Times > j.(SortDemoStruct).Times
			}
		}
	} else {
		fmt.Println("myFunc init error")
		os.Exit(1)
	}
}

// 发送数据到服务器
func (this *QuickSortStruct) SendDataToServe(conn net.Conn) {
	if conn == nil {
		fmt.Println("client is empty")
		return
	}
	if this.IsFile && this.SendFilePath == "" {
		fmt.Println("send file is empty")
		os.Exit(404)
	}
	// sent to client
	byte0 := IntToBytes(0)
	byte1 := IntToBytes(1)

	var engine []byte
	if this.IsFile {
		engine = byte1
	} else {
		engine = byte0
	}
	var sortByte []byte
	if this.IsAsc {
		sortByte = byte1
	} else {
		sortByte = IntToBytes(2)
	}
	start := append(append(engine, byte0...), byte0...)
	conn.Write(start)
	// 获取数据长度
	length := 0
	if this.IsFile {
		length = GetFileLineNumber(this.SendFilePath)
		open, err := os.Open(this.SendFilePath)
		ManageError(err)
		reader := bufio.NewReader(open)
		// 从文件中读取数据 并转成对应的格式
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			tmpStr := string(line)
			if tmpStr == "" {
				continue
			}
			toByte := make([]byte, 0)
			if this.DataType == "int" {
				tmp, _ := strconv.Atoi(tmpStr)
				toByte = this.DataFormatToByte(tmp)
			} else if "float64" == this.DataType {
				tmp, _ := strconv.ParseFloat(tmpStr, 64)
				toByte = this.DataFormatToByte(tmp)
			} else if "string" == this.DataType {
				toByte = this.DataFormatToByte(tmpStr)
			} else if "struct" == this.DataType {
				split := strings.Split(tmpStr, " # ")
				s := new(SortDemoStruct)
				s.Info = split[0]
				s.Times, _ = strconv.Atoi(split[1])
				toByte = this.DataFormatToByte(s)
			} else {
				fmt.Println("error data type")
				os.Exit(404)
			}
			//toByte := u.DataFormatToByte(tmp)
			conn.Write(append(engine, toByte...))
		}
	} else {
		length = len(this.Data)
		for i := 0; i < length; i++ {
			toByte := this.DataFormatToByte(this.Data[i])
			conn.Write(append(engine, toByte...))
		}
	}
	end := append(append(engine, byte0...), sortByte...)
	conn.Write(end)
}

// 从服务器接收数据
func (this *QuickSortStruct) ReceiveDataFromServe(conn net.Conn) {
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
			d1 := BytesToInt(byt[:8])
			d2 := BytesToInt(byt[8:16])
			d3 := BytesToInt(byt[16:])
			// 存储文件排序
			if d1 == 1 {
				this.IsFile = true
				// 获取文件存储位置
				if this.ReceiveFilePath == "" {
					this.GetReceiveDataSavePath(false)
				}
			}
			// 表示开始数据传输
			if d2 == 0 && d3 == 0 {
				fmt.Println("start receive data")
				this.Data = make([]interface{}, 0, 0)
			}
			// 接收对应的数据类型
			if d2 == 1 { // 整数
				this.DataType = "int"
				this.PutEveryReceiveToData(d3)
			} else if d2 == 2 { // 浮点数
				this.DataType = "float64"
				this.PutEveryReceiveToData(BytesToFloat64(byt[16:]))
			} else if d2 == 3 { // 字符串
				strByt := make([]byte, d3, d3)
				read, _ := conn.Read(strByt)
				if read == d3 {
					this.DataType = "string"
					this.PutEveryReceiveToData(string(strByt))
				}
			} else if d2 == 4 { // 结构体
				jsonByt := make([]byte, d3, d3)
				read, _ := conn.Read(jsonByt)
				if read == d3 {
					this.DataType = "struct"
					this.PutEveryReceiveToData(jsonByt)
				}
			}
			// 数据接收结束
			if d2 == 0 && (d3 == 1 || d3 == 2) {
				fmt.Println("finish receive data")
				if this.IsFile {
					this.ReceiveFileWrite.Flush()
					fmt.Println(this.ReceiveFilePath)
				} else {
					fmt.Println(this.Data)
				}
			}
		}
	}
}

// 已经排序好的文件合并成一个文件
func (this *QuickSortStruct) MergeFileListAsOne(fileList []string, savePath string) string {
	l := list.New()
	for i := 0; i < len(fileList); i++ {
		l.PushBack(fileList[i])
	}
	fmt.Println("file total is:", l.Len())
	// 归并合并文件直到栈内剩余一个文件
	i := 0
	for l.Len() != 1 {
		file1 := l.Front()
		l.Remove(file1)

		file2 := l.Front()
		l.Remove(file2)

		if file1 != nil && file2 != nil {
			f1 := file1.Value.(string)
			f2 := file2.Value.(string)
			one := this.MergeTwoFileAsOne(f1, f2, savePath, i)
			if one == "" {
				fmt.Println("file merge error")
				os.Exit(404)
			}
			l.PushBack(one)
			i++
		} else if file1 != nil {
			f1 := file1.Value.(string)
			l.PushBack(f1)
		} else if file2 != nil {
			f2 := file2.Value.(string)
			l.PushBack(f2)
		} else {
			break
		}
	}
	back := l.Back()
	l.Remove(back)
	return back.Value.(string)
}

// 两个文件归并
func (this *QuickSortStruct) MergeTwoFileAsOne(file1, file2, path string, i int) string {
	line1 := GetFileLineNumber(file1)
	line2 := GetFileLineNumber(file2)

	nano := time.Now().UnixNano()
	path += strconv.Itoa(int(nano)) + "_" + strconv.Itoa(i) + ".txt"
	create, err := os.Create(path)
	defer create.Close()
	if err != nil {
		fmt.Println("file create err")
		return ""
	}
	writer := bufio.NewWriter(create)

	f1, err := os.Open(file1)
	defer f1.Close()
	if err != nil {
		fmt.Println("file open err")
		return ""
	}
	r1 := bufio.NewReader(f1)

	f2, err := os.Open(file2)
	defer f2.Close()
	if err != nil {
		fmt.Println("file open err")
		return ""
	}
	r2 := bufio.NewReader(f2)

	i, j := 0, 0
	tmp1, _, _ := r1.ReadLine()
	tmp2, _, _ := r2.ReadLine()
	str1 := string(tmp1)
	str2 := string(tmp2)

	for i < line1 && j < line2 {
		// 根据类型还原最终数据 并进行大小比较 0：相等 1：d1>d2 2：d1<d2
		res := 0 // 0：相等 1：d1>d2 2：d1<d2
		if this.DataType == "int" {
			d1, _ := strconv.Atoi(str1)
			d2, _ := strconv.Atoi(str2)
			res = this.DataCompare(d1, d2)
		} else if "float64" == this.DataType {
			d1, _ := strconv.ParseFloat(str1, 64)
			d2, _ := strconv.ParseFloat(str2, 64)
			res = this.DataCompare(d1, d2)
		} else if "string" == this.DataType {
			d1 := str1
			d2 := str2
			res = this.DataCompare(d1, d2)
		} else if "struct" == this.DataType {
			split1 := strings.Split(str1, " # ")
			d1 := new(SortDemoStruct)
			d1.Info = split1[0]
			d1.Times, _ = strconv.Atoi(split1[1])

			split2 := strings.Split(str2, " # ")
			d2 := new(SortDemoStruct)
			d2.Info = split2[0]
			d2.Times, _ = strconv.Atoi(split2[1])
			res = this.DataCompare(d1, d2)
		} else {
			fmt.Println("error data type")
			os.Exit(404)
		}

		// 大小比较结果 0：相等 1：d1>d2 2：d1<d2
		if 1 == res {
			if this.IsAsc {
				fmt.Fprintln(writer, str1)
				tmp1, _, _ = r1.ReadLine()
				str1 = string(tmp1)
				i++
			} else {
				fmt.Fprintln(writer, str2)
				tmp2, _, _ = r2.ReadLine()
				str2 = string(tmp2)
				j++
			}
		} else if 2 == res {
			if this.IsAsc {
				fmt.Fprintln(writer, str2)
				tmp2, _, _ = r2.ReadLine()
				str2 = string(tmp2)
				j++
			} else {
				fmt.Fprintln(writer, str1)
				tmp1, _, _ = r1.ReadLine()
				str1 = string(tmp1)
				i++
			}
		} else {
			fmt.Fprintln(writer, str1)
			tmp1, _, _ = r1.ReadLine()
			str1 = string(tmp1)
			i++

			fmt.Fprintln(writer, str2)
			tmp2, _, _ = r2.ReadLine()
			str2 = string(tmp2)
			j++
		}
	}
	for i < line1 {
		line, _, _ := r1.ReadLine()
		str := string(line)
		if str != "" {
			fmt.Fprintln(writer, str)
		}
		i++
	}
	for j < line2 {
		line, _, _ := r2.ReadLine()
		str := string(line)
		if str != "" {
			fmt.Fprintln(writer, str)
		}
		j++
	}
	writer.Flush()
	create.Close()
	f1.Close()
	f2.Close()
	return path
}

// 大小比较 0：相等 1：d1>d2 2：d1<d2
func (this *QuickSortStruct) DataCompare(d1, d2 interface{}) int {
	res := 0
	if this.DataType == "int" {
		if d1.(int) == d2.(int) {
			return 0
		} else if d1.(int) > d2.(int) {
			return 1
		} else {
			return 2
		}
	} else if "float64" == this.DataType {
		if d1.(float64) == d2.(float64) {
			return 0
		} else if d1.(float64) > d2.(float64) {
			return 1
		} else {
			return 2
		}
	} else if "string" == this.DataType {
		if d1.(string) == d2.(string) {
			return 0
		} else if d1.(string) > d2.(string) {
			return 1
		} else {
			return 2
		}
	} else if "struct" == this.DataType {
		if d1.(SortDemoStruct).Times == d2.(SortDemoStruct).Times {
			return 0
		} else if d1.(SortDemoStruct).Times > d2.(SortDemoStruct).Times {
			return 1
		} else {
			return 2
		}
	} else {
		fmt.Println("error data type")
		os.Exit(404)
	}
	return res
}
