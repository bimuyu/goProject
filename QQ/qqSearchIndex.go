package QQ

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type QQAccount struct {
	QQNumber int
	QQPwd    string
}

// 清洗不合格qq数据
func ClearQQData(source, save, errorPath string) {
	open, _ := os.Open(source)
	defer open.Close()
	reader := bufio.NewReader(open)

	create, _ := os.Create(save)
	defer create.Close()
	right1 := bufio.NewWriter(create)

	wrongPath, _ := os.Create(errorPath)
	defer wrongPath.Close()
	wrong := bufio.NewWriter(wrongPath)

	reg1 := `^[1-9][0-9]{4,11}$`
	reg2 := `^\S{6,25}$`
	regx1 := regexp.MustCompile(reg1)
	regx2 := regexp.MustCompile(reg2)
	i := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		split := strings.Split(tmp, "----")
		if len(split) == 2 {
			if regx1.Match([]byte(split[0])) && regx2.Match([]byte(split[1])) {
				fmt.Fprintln(right1, tmp)
			} else {
				fmt.Fprintln(wrong, tmp)
			}
		} else {
			fmt.Fprintln(wrong, tmp)
		}
		i++
		if i%5000000 == 0 {
			fmt.Println(i)
		}
	}
	right1.Flush()
	wrong.Flush()
	fmt.Println("所有qq号码提取并写入完成")
}

// 清洗不合格qq数据
func QQDataOrderByQQNumber(source, save string) {
	//number := Utils.GetFileLineNumber(source)
	//fmt.Println(number)
	//os.Exit(1)
	number := 156730744 //qq
	fmt.Println("开始时间：", time.Now().Format("2006-01-02 15:04:05"))
	pwd := make([]QQAccount, number, number)
	fmt.Println("数组大小：", len(pwd))

	// 读取文件到结构体
	pwd = readQQDataToStruct(pwd, source)
	fmt.Println("写入内存完成时间：", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("数组大小：", len(pwd))

	// 结构体排序 从多到少
	pwd = qqDataStructSort(pwd)
	fmt.Println("排序完成时间：", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("数组大小：", len(pwd))

	// 排序数据重新写入文件
	qqAccountStructSaveToFile(pwd, save)
	fmt.Println("写入文件完成时间：", time.Now().Format("2006-01-02 15:04:05"))
}

// 读取文件到结构体
func readQQDataToStruct(pwd []QQAccount, path string) []QQAccount {
	open, _ := os.Open(path)
	reader := bufio.NewReader(open)
	i := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		if tmp != "" {
			split := strings.Split(tmp, "----")
			if len(split) == 2 {
				atoi, _ := strconv.Atoi(split[0])
				pwd[i].QQNumber = atoi
				pwd[i].QQPwd = split[1]
				i++
			}
		}
	}
	open.Close()
	return pwd
}

func qqDataStructSort(pwd []QQAccount) []QQAccount {
	if len(pwd) < 100 {
		return qqDataStructInsertSort(pwd)
	}
	qqAccountStructQuickSort(pwd, 0, len(pwd)-1)
	return pwd
}

func qqDataStructInsertSort(pwd []QQAccount) []QQAccount {
	if len(pwd) <= 1 {
		return pwd
	}
	for i := 1; i < len(pwd); i++ {
		index := findQQAccountStructInsertIndex(pwd, 0, i-1, i)
		if i != index {
			for j := i; j > index; j-- {
				pwd[j], pwd[j-1] = pwd[j-1], pwd[j]
			}
		}
	}
	return pwd
}

func findQQAccountStructInsertIndex(pwd []QQAccount, start, end, index int) int {
	if start >= end {
		if pwd[start].QQNumber < pwd[index].QQNumber {
			return start
		} else {
			return index
		}
	}
	mid := (start + end) / 2
	if pwd[mid].QQNumber == pwd[index].QQNumber {
		return mid
	} else if pwd[mid].QQNumber < pwd[index].QQNumber {
		return findQQAccountStructInsertIndex(pwd, start, mid, index)
	} else {
		return findQQAccountStructInsertIndex(pwd, mid+1, end, index)
	}
}

func qqAccountStructQuickSort(pwd []QQAccount, start, end int) {
	if end-start <= 100 {
		qqAccountStructQuickSortPart(pwd, start, end)
	} else {
		swapQQAccount(pwd, start, rand.Int()%(end-start+1)+start)
		tmp := pwd[start] // 中间锚点
		lt := start       // pwd[start+1,lt] >tmp.times lt++
		gt := end + 1     // pwd[gt,end] <tmp.times gt--
		i := start + 1    // pwd[lt+1,i] == tmp.times i++
		for i < gt {
			if pwd[i].QQNumber > tmp.QQNumber {
				swapQQAccount(pwd, i, lt+1)
				lt++
				i++
			} else if pwd[i].QQNumber < tmp.QQNumber {
				swapQQAccount(pwd, i, gt-1)
				gt--
			} else {
				i++
			}
		}
		swapQQAccount(pwd, start, lt)
		qqAccountStructQuickSort(pwd, start, lt-1)
		qqAccountStructQuickSort(pwd, gt, end)
	}
}

func qqAccountStructQuickSortPart(pwd []QQAccount, start, end int) []QQAccount {
	if end-start <= 1 {
		if pwd[start].QQNumber < pwd[end].QQNumber {
			swapQQAccount(pwd, start, end)
		}
	}
	for i := start + 1; i <= end; i++ {
		index := findQQAccountStructInsertIndex(pwd, start, i-1, i)
		if i != index {
			for j := i; j > index; j-- {
				pwd[j], pwd[j-1] = pwd[j-1], pwd[j]
			}
		}
	}
	return pwd
}

func swapQQAccount(pwd []QQAccount, i, j int) {
	pwd[i], pwd[j] = pwd[j], pwd[i]
}

func qqAccountStructSaveToFile(pwd []QQAccount, save string) {
	create, _ := os.Create(save)
	writer := bufio.NewWriter(create)
	writer = bufio.NewWriterSize(writer, 1024*100)
	//n := len(pwd) / 2
	//pwd = pwd[:1000000]
	for i := 0; i < len(pwd); i++ {
		fmt.Fprintln(writer, strconv.Itoa(pwd[i].QQNumber)+"----"+pwd[i].QQPwd)
	}
	writer.Flush()
	create.Close()
}

// 排序好的qq账号数据读取的内存
func ReadQQAccountSortDataToMemory(source string) []string {
	//number := Utils.GetFileLineNumber(source)
	//fmt.Println(number)
	//os.Exit(1)
	number := 156730744
	arr := make([]string, number, number)
	file, _ := os.Open(source)
	reader := bufio.NewReader(file)
	i := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		if tmp != "" {
			arr[i] = tmp
		}
		i++
		if i%10000000 == 0 {
			fmt.Println(i)
		}
	}
	defer file.Close()
	return arr
}

// 按照结构体读入文件
func ReadQQAccountSortDataToMemoryAsStruct(source string) []QQAccount {
	//number := Utils.GetFileLineNumber(source)
	//fmt.Println(number)
	//os.Exit(1)
	number := 156730744
	arr := make([]QQAccount, number, number)
	file, _ := os.Open(source)
	reader := bufio.NewReader(file)
	i := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		if tmp != "" {
			split := strings.Split(tmp, "----")
			arr[i].QQPwd = split[1]
			arr[i].QQNumber, _ = strconv.Atoi(split[0])
		}
		i++
		if i%10000000 == 0 {
			fmt.Println(i)
		}
	}
	defer file.Close()
	return arr
}

// 二分查找
func BinarySearchQQPwdFromStruct(arr []QQAccount, qq int) int {
	left := 0
	right := len(arr) - 1
	for left <= right {
		mid := (left + right) / 2
		if arr[mid].QQNumber == qq {
			return mid
		} else if arr[mid].QQNumber > qq {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// 生成文件随机访问索引数组 把根据每行长度 叠加生成每行数据的位置
func MakeQQAccountSearchIndexArray(sourcePath string) []int {
	//number := Utils.GetFileLineNumber(source)
	//fmt.Println(number)
	//os.Exit(1)
	number := 156730744
	index := make([]int, number+1, number+1)

	file, _ := os.Open(sourcePath)
	defer file.Close()
	reader := bufio.NewReader(file)
	index[0] = 0
	i := 1
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		// 加上最后的换行符
		index[i] = len(tmp) + 1
		i++
		if i%10000000 == 0 {
			fmt.Println(i)
		}
	}
	for j := 0; j < len(index)-1; j++ {
		index[j+1] += index[j]
	}
	return index
}

// 根据行号访问文件
func SearchQQAccountByLineNumber(arr []int, lineNumber int, file *os.File) {
	file.Seek(0, 0)
	linePos := arr[lineNumber]
	file.Seek(int64(linePos), 0)
	bytes := make([]byte, 30, 30)
	file.Read(bytes)
	pos := 0
	for i := 0; i < len(bytes); i++ {
		if i >= 14 && bytes[i] == '\n' {
			pos = i
		}
	}
	fmt.Println("第", lineNumber, "行的数据是：", string(bytes[:pos]))

	fmt.Println("请再次输入行号")
}

// 快速搜索qq账号的密码
func SearchQQAccountPwdUseSearchIndexArray(arr []int, qq int, file *os.File) {
	left := 0
	right := len(arr) - 1

	// 快速获取qq在文件中的行号
	str := strconv.Itoa(qq)
	fmt.Println("need find data is:", str)
	for left <= right {
		mid := (left + right) / 2
		file.Seek(int64(arr[mid]), 0)
		bytes := make([]byte, 34, 34)
		n, _ := file.Read(bytes)
		pos := 0
		for j := 0; j < n-1; j++ {
			if j >= 14 && bytes[j] == '\n' {
				pos = j
				break
			}
		}
		tmp := string(bytes[:pos])
		split := strings.Split(tmp, "----")
		midQQ, _ := strconv.Atoi(split[0])
		if midQQ == qq {
			fmt.Println("data is:", tmp)
			break
		} else if midQQ > qq {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	fmt.Println("搜索完毕，请输入qq号")
	return
}

// qq账号索引文件保存到文件
func PutQQAccountSearchIndexToFile(source, save string) {
	array := MakeQQAccountSearchIndexArray(source)
	fmt.Println("start put to file")
	file, _ := os.Create(save)
	defer file.Close()
	writer := bufio.NewWriter(file)
	for i := 0; i < len(array); i++ {
		len := array[i]
		bytes := make([]byte, 4, 4)
		binary.BigEndian.PutUint32(bytes, uint32(len))
		writer.Write(bytes)
		if i%10000000 == 0 {
			fmt.Println(i)
		}
	}
	writer.Flush()
}

// 使用查询索引文件 根据行号快速查询
func QuickFindQQAccountDataUseSearchIndexForLineNumber(sourcePath, indexPath string) string {
	file, _ := os.Open(sourcePath)
	fmt.Println("请输入行号")
	index, _ := os.Open(indexPath)
	defer func() {
		file.Close()
		index.Close()
	}()

	file.Seek(0, 0)
	index.Seek(0, 0)

	for {
		var lineNumber int
		fmt.Scanf("%d", &lineNumber)

		pos := lineNumber * 4
		bytes := make([]byte, 4)
		index.Seek(int64(pos), 0)
		index.Read(bytes)
		linePos := binary.BigEndian.Uint32(bytes)
		file.Seek(int64(linePos), 0)
		lineByte := make([]byte, 34, 34)
		n, _ := file.Read(lineByte)
		endPos := 0
		for i := 0; i < n; i++ {
			if i >= 14 && lineByte[i] == '\n' {
				endPos = i
				break
			}
		}
		return string(lineByte[:endPos])
	}
}

// 使用查询索引文件 快速查询qq密码
func QuickFindQQAccountPwdUseSearchIndex(qq int, source, index *os.File) {
	end := 156730744
	start := 0
	for start <= end {
		mid := (start + end) / 2
		tmp, pwd := findLineInfoByLineNumber(mid, source, index)
		if tmp == qq {
			fmt.Println("qq:", qq, ",pwd:", pwd)
			return
		} else if tmp > qq {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	fmt.Println("not fount")
	fmt.Println("请输QQ号")
}
func findLineInfoByLineNumber(lineNumber int, source, index *os.File) (int, string) {
	index.Seek(int64(lineNumber*4), 0)
	lineLengthByte := make([]byte, 4)
	index.Read(lineLengthByte)

	linePos := binary.BigEndian.Uint32(lineLengthByte)
	source.Seek(int64(linePos), 0)
	lineByte := make([]byte, 34)
	n, _ := source.Read(lineByte)
	endPos := 0
	for i := 0; i < n; i++ {
		if i > 14 && lineByte[i] == '\n' {
			endPos = i
			break
		}
	}
	line := string(lineByte[:endPos])

	split := strings.Split(line, "----")
	qq, _ := strconv.Atoi(split[0])
	return qq, split[1]
}
