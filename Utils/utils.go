package Utils

import (
	"bufio"
	"bytes"
	"container/list"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 二分查找
func BinSearch(arr []int, pot int) int {
	left := 0
	right := len(arr) - 1
	for left <= right {
		mid := (left + right) / 2
		fmt.Println(mid)
		if arr[mid] < pot {
			left = mid + 1
		} else if arr[mid] > pot {
			right = mid - 1
		} else {
			return mid
		}
	}
	return -1
}

// mysql链接
func MysqlConnect(host, port, password, user, dbName, charset string) bool {
	// username:password@tcp(ip:port)/database?charset=utf8
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=" + charset
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("connect error", err)
		return false
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return false
	}
	res, err := db.Query("select count(*) as number from area")

	for res.Next() {
		var number int
		err := res.Scan(&number)
		if err != nil {
			fmt.Println("get data error", err)
			return false
		}
		fmt.Println(number)
	}
	return true
}

// 清洗数据到新文件
func DataClearToNewFile(sourcePath, savePath string) {
	source, err := os.Open(sourcePath)
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	reader := bufio.NewReader(source)
	create, err := os.Create(savePath)
	if err != nil {
		fmt.Println("file create error:", err)
		return
	}
	writer := bufio.NewWriter(create)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		split := strings.Split(tmp, "|")
		if 5 == len(split) {
			inn := split[len(split)-1]
			writer.WriteString(inn + "\n")
		}
	}
	writer.Flush()
}

// 获取文件的行数
func GetFileLineNumber(path string) int {
	open, err := os.Open(path)
	if err != nil {
		fmt.Println("file open error:", err)
		return -1
	}
	reader := bufio.NewReader(open)
	n := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if string(line) != "" {
			n++
		}
	}
	open.Close()
	return n
}

// 快速排序
func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	tmp := arr[0]
	var equal, big, less []int
	equal = append(equal, tmp)
	for i := 1; i < len(arr); i++ {
		if arr[i] == tmp {
			equal = append(equal, tmp)
		} else if arr[i] > tmp {
			big = append(big, arr[i])
		} else {
			less = append(less, arr[i])
		}
	}
	less, big = QuickSort(less), QuickSort(big)
	return append(append(less, equal...), big...)
}

// 生产随机数字数组
func MakeArrFunc(number int) []int {
	var arr []int
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < number; i++ {
		intn := source.Intn(1000)
		arr = append(arr, intn)
	}
	return arr
}

// 把数据读入数组并排序
func ReadFileToArrayAndSortDemo(sourcePath string, arr []int) []int {
	open, err := os.Open(sourcePath)
	defer open.Close()
	if err != nil {
		fmt.Println("file open error")
		return nil
	}
	reader := bufio.NewReader(open)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmpStr := string(line)
		tmpInt, _ := strconv.Atoi(tmpStr)
		arr = append(arr, tmpInt)
	}
	fmt.Println("数据读入数组完成")
	return QuickSort(arr)
}

// 把数据读入数组并排序
func SaveArrayToFile(savePath string, arr []int) {
	create, err := os.Create(savePath)
	defer create.Close()
	if err != nil {
		fmt.Println("file create err")
		return
	}
	writer := bufio.NewWriter(create)
	for i := 0; i < len(arr); i++ {
		tmp := strconv.Itoa(arr[i])
		writer.WriteString(tmp + "\n")
	}
	writer.Flush()
}

// 数据去重并计算次数
func DataUniqueAndCount(sourcePath string, savePath string) {
	open, err := os.Open(sourcePath)
	if err != nil {
		fmt.Println("file open err")
		return
	}
	defer open.Close()

	create, err := os.Create(savePath)
	if err != nil {
		fmt.Println("file create err")
		return
	}
	defer create.Close()

	reader := bufio.NewReader(open)
	writer := bufio.NewWriter(create)

	tmp := ""
	times := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmpStr := string(line)
		if tmp == "" || tmp != tmpStr {
			if tmp != "" {
				writer.WriteString(tmp + " # " + strconv.Itoa(times) + "\n")
			}
			tmp = tmpStr
			times = 0
		}
		times++
	}
	writer.Flush()
}

type DemoStruct struct {
	info  string
	times int
}

// 数据写入结构体
func DataSortByTimes(sourcePath string) []DemoStruct {
	open, err := os.Open(sourcePath)
	if err != nil {
		fmt.Println("file open err")
		return nil
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	var demoStruct []DemoStruct
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmpStr := string(line)
		split := strings.Split(tmpStr, " # ")
		var ds DemoStruct
		tmpInt, _ := strconv.Atoi(split[1])
		ds.info = split[0]
		ds.times = tmpInt
		demoStruct = append(demoStruct, ds)
	}
	return demoStruct
}

// 结构体快速排序
func DemoStructInfoQuickSort(demo []DemoStruct) []DemoStruct {
	if len(demo) <= 1 {
		return demo
	}
	var equal, less, more []DemoStruct
	equal = append(equal, demo[0])
	tmp := demo[0]
	for i := 1; i < len(demo); i++ {
		ds := demo[i]
		if ds.times > tmp.times {
			more = append(more, ds)
		} else if ds.times < tmp.times {
			less = append(less, ds)
		} else {
			equal = append(equal, ds)
		}
	}
	less, more = DemoStructInfoQuickSort(less), DemoStructInfoQuickSort(more)
	return append(append(more, equal...), less...)
}

// 排序后的结构体写入文件
func SaveDemoStructToFile(demo []DemoStruct, savePath string) {
	create, err := os.Create(savePath)
	if err != nil {
		fmt.Println("file create err")
		return
	}
	defer create.Close()
	writer := bufio.NewWriter(create)
	for i := 0; i < len(demo); i++ {
		ds := demo[i]
		writer.WriteString(ds.info + " # " + strconv.Itoa(ds.times) + "\n")
	}
	writer.Flush()
}

// 简单插入排序
func SimpleInsertSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	for i := 1; i < len(arr); i++ {
		insertVal := arr[i]
		insertIndex := i - 1
		for insertIndex >= 0 && insertVal < arr[insertIndex] {
			arr[insertIndex+1] = arr[insertIndex]
			insertIndex--
		}
		arr[insertIndex+1] = insertVal
	}
	return arr
}

func InsertSortDemo() {
	arr := []int{11, 21, 6, 3, 9, 23, 5, 45, 45, 232, 54}
	insertVal := arr[2]
	insertIndex := 2 - 1
	for insertIndex >= 0 && arr[insertIndex] > insertVal {
		arr[insertIndex+1] = arr[insertIndex]
		insertIndex--
	}
	arr[insertIndex+1] = insertVal
	fmt.Println(arr)
}

// 二分插入排序
func BinInsertSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	length := len(arr)
	for i := 1; i < length; i++ {
		index := findInsertIndex(arr, 0, i-1, i)
		if index != i {
			for j := i; j > index; j-- {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			}
		}
	}
	return arr
}

// 获取数据要插入的位置
func findInsertIndex(arr []int, start, end, current int) int {
	if start >= end {
		if arr[start] > arr[current] {
			return start
		} else {
			return current
		}
	}
	mid := (start + end) / 2
	if arr[mid] < arr[current] {
		return findInsertIndex(arr, mid+1, end, current)
	} else if arr[mid] == arr[current] {
		return mid
	} else {
		return findInsertIndex(arr, start, mid, current)
	}
}

// 优化快速排序
func OptimizeQuickSort(arr []int) []int {
	// todo 数据量小于100时插入排序效率更高
	if len(arr) < 20 {
		return BinInsertSort(arr)
	}
	optimizeQuickSort(arr, 0, len(arr)-1)
	return arr
}

// 交换
func Swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// 优化快速排序
func optimizeQuickSort(arr []int, left, right int) {
	if right-left < 10 { //截取数组的一段数据 并进行插入排序
		sortPartOfArray(arr, left, right)
	} else {
		Swap(arr, left, rand.Int()%(right-left+1)+left)
		tmp := arr[left] // 快速排序 锚点值
		lt := left       // 定义使得arr[left+1,lt] 范围内的数据都<tmp lt++
		gt := right + 1  // 定义使得arr[gt,right] 范围内的数据都>tmp gt--
		i := left + 1    // 定义使得arr[lt+1,i] 范围内的数据都=tmp i++
		for i < gt {
			if arr[i] > tmp {
				Swap(arr, i, gt-1)
				gt--
			} else if arr[i] < tmp {
				Swap(arr, i, lt+1)
				i++
				lt++
			} else {
				i++
			}
		}
		//fmt.Println("left:", left)
		//fmt.Println("right:", right)
		//fmt.Println("lt:", lt)
		//fmt.Println("gt:", gt)
		//fmt.Println("arr:", arr)
		Swap(arr, left, lt)
		//fmt.Println("arr:", arr)
		//os.Exit(0)
		if lt-1 > left {
			optimizeQuickSort(arr, left, lt-1)
		}
		optimizeQuickSort(arr, gt, right)
	}
}

// 截取数组的一段数据 并进行插入排序
func sortPartOfArray(arr []int, start, end int) []int {
	if end-start <= 1 {
		// 判断大小 根据排序要求交换位置
		if arr[start] > arr[end] {
			Swap(arr, start, end)
		}
		return arr
	}
	for i := start + 1; i <= end; i++ {
		index := findInsertIndex(arr, start, i-1, i)
		if i != index {
			for j := i; j > index; j-- {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			}
		}
	}
	return arr
}

// 遍历文件夹
func TraverseDir(path string) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("file open err")
		return
	}
	for _, v := range dir {
		if v.IsDir() {
			fmt.Println(path + "/" + v.Name())
			TraverseDir(path + "/" + v.Name())
		} else {
			if strings.Index(v.Name(), ".") != 0 {
				fmt.Println(path+"/"+v.Name(), v.Size())
			}
		}
	}
}

// 均等切割数据
func EqualSplitData(num, n int) []int {
	if num < n {
		return nil
	}
	arr := make([]int, n, n)
	if num%n == 0 {
		tmp := num / n
		for i := 0; i < n; i++ {
			arr[i] = tmp
		}
	} else {
		tmp := (num - num%n) / (n - 1)
		for i := 0; i < n-1; i++ {
			arr[i] = tmp
			num -= tmp
		}
		arr[n-1] = num
	}
	return arr
}

// 把文件均等切割
func SliceFileToSmall(sourcePath string, savePath string, n int) []string {
	number := GetFileLineNumber(sourcePath)
	data := EqualSplitData(number, n)
	open, err := os.Open(sourcePath)
	if err != nil {
		fmt.Println("file open error")
		return nil
	}
	defer open.Close()
	res := make([]string, n, n)
	reader := bufio.NewReader(open)
	for i := 0; i < len(data); i++ {
		tmpPath := savePath + strconv.Itoa(i) + ".txt"
		res[i] = tmpPath
		create, err := os.Create(tmpPath)
		if err != nil {
			fmt.Println("file create error")
			return nil
		}
		writer := bufio.NewWriter(create)
		for j := 0; j < data[i]; j++ {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			tmp := string(line)
			fmt.Fprintln(writer, tmp)
		}
		writer.Flush()
		create.Close()
	}
	open.Close()
	return res
}

// 数组归并排序
func ArrayMergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	var wg sync.WaitGroup
	wg.Add(2)
	var left, right []int
	go func() {
		left = ArrayMergeSort(arr[:mid])
		wg.Done()
	}()
	go func() {
		right = ArrayMergeSort(arr[mid:])
		wg.Done()
	}()
	wg.Wait()
	return TwoArrayMerge(left, right)
}
func TwoArrayMerge(left, right []int) []int {
	result := []int{}
	i, j := 0, 0
	n1, n2 := len(left), len(right)

	for i < n1 && j < n2 {
		if left[i] < right[j] {
			result = append(result, right[j])
			j++
		} else if left[i] > right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
			result = append(result, left[i])
			i++
		}
	}
	if i != n1 {
		result = append(result, left[i:]...)
	}
	if j != n2 {
		result = append(result, right[j:]...)
	}
	return result
}

// 整型转字节
func IntToBytes(n int) []byte {
	// 数据类型转换
	i := int64(n)
	// 字节集合
	buffer := bytes.NewBuffer([]byte{})
	// 按照二进制写入集合
	binary.Write(buffer, binary.BigEndian, i)
	// 返回字节数组
	return buffer.Bytes()
}

// 字节转整数
func BytesToInt(byt []byte) int {
	// 写入二进制集合
	buffer := bytes.NewBuffer(byt)
	var data int64
	// 解码到定义的数据
	binary.Read(buffer, binary.BigEndian, &data)
	return int(data)
}

func Float32ToBytes(n float32) []byte {
	bits := math.Float32bits(n)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
func Float64ToBytes(n float64) []byte {
	bits := math.Float64bits(n)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
func BytesToFloat32(byt []byte) float32 {
	u := binary.LittleEndian.Uint32(byt)
	return math.Float32frombits(u)
}
func BytesToFloat64(byt []byte) float64 {
	u := binary.LittleEndian.Uint64(byt)
	return math.Float64frombits(u)
}

type SortDemoStruct struct {
	Info  string
	Times int
}

func JsonDemo() {
	u := new(SortDemoStruct)
	u.Info = "this is test"
	u.Times = 10
	bytes, _ := json.Marshal(u)
	fmt.Println(bytes)
	tmp := new(SortDemoStruct)
	json.Unmarshal(bytes, tmp)
	fmt.Println(tmp.Info)
	fmt.Println(tmp.Times)
}

// 已经排序好的文件合并成一个文件
func MergeFileListAsOne(fileList []string) string {
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
			one := MergeTwoFileAsOne(f1, f2, i)
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
func MergeTwoFileAsOne(file1, file2 string, i int) string {
	path := "/Users/magic/web/golang/source/csdn_split/"
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
		if str1 < str2 {
			if str2 != "" {
				writer.WriteString(str2 + "\n")
				tmp2, _, _ = r2.ReadLine()
				str2 = string(tmp2)
			}
			j++
		} else if str1 > str2 {
			if str1 != "" {
				writer.WriteString(str1 + "\n")
				tmp1, _, _ = r1.ReadLine()
				str1 = string(tmp1)
			}
			i++
		} else {
			if str1 != "" {
				writer.WriteString(str1 + "\n")
				tmp1, _, _ = r1.ReadLine()
				str1 = string(tmp1)
			}
			i++
			if str2 != "" {
				writer.WriteString(str2 + "\n")
				tmp2, _, _ = r2.ReadLine()
				str2 = string(tmp2)
			}
			j++
		}
	}
	for i < line1 {
		line, _, _ := r1.ReadLine()
		str := string(line)
		if str != "" {
			writer.WriteString(str + "\n")
		}
		i++
	}
	for j < line2 {
		line, _, _ := r2.ReadLine()
		str := string(line)
		if str != "" {
			writer.WriteString(str + "\n")
		}
		j++
	}
	writer.Flush()
	create.Close()
	f1.Close()
	f2.Close()
	return path
}