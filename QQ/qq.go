package QQ

import (
	"bigData/Utils"
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetQQPwd(source, save string) {
	open, err := os.Open(source)
	if err != nil {
		fmt.Println("file open err")
		return
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	create, err := os.Create(save)
	if err != nil {
		fmt.Println("file create err")
		return
	}
	defer create.Close()
	writer := bufio.NewWriter(create)
	reg1 := `^[1-9][0-9]{4,11}$`
	reg2 := `^.{6,25}$`
	regx1 := regexp.MustCompile(reg1)
	regx2 := regexp.MustCompile(reg2)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		split := strings.Split(tmp, "----")
		if len(split) == 2 {
			//if regx1.MatchString(split[1]) && regx2.MatchString(split[1]) {
			if regx1.Match([]byte(split[0])) && regx2.Match([]byte(split[1])) {
				//writer.WriteString(split[1] + "\n")
				fmt.Fprintln(writer, split[1])
			}
		}
	}
	writer.Flush()
	fmt.Println("所有qq密码提取并写入完成")
}

func QQPwdQuickSort(source, save string) {
	// 获取密码行数
	number := Utils.GetFileLineNumber(source)
	fmt.Println(number)
	//os.Exit(1)
	//number := 6428632
	//number := 156783175
	fmt.Println("总数据量：", number)
	// 读取文件到内存
	open, err := os.Open(source)
	defer open.Close()
	if err != nil {
		fmt.Println("file open error")
		return
	}
	reader := bufio.NewReader(open)
	fmt.Println("开始读取文件到内存", time.Now().Unix())
	tmp := make([]string, number+1, number+1)
	fmt.Println("数组大小", len(tmp))
	for i := 0; i < number; i++ {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		str := string(line)
		if str != "" {
			tmp[i] = str
		}
	}
	fmt.Println("数组大小", len(tmp))
	// 回收内存资源
	open.Close()
	runtime.GC()
	debug.FreeOSMemory()

	fmt.Println("文件开始排序", time.Now().Unix())
	// 排序
	sort := optimizeQuickSort(tmp)
	// 回收内存资源
	tmp = nil
	runtime.GC()
	debug.FreeOSMemory()

	fmt.Println("文件排序完成", time.Now().Unix())
	// 排序数据写入文件
	create, err := os.Create(save)
	defer create.Close()
	if err != nil {
		fmt.Println("file create error")
		return
	}
	writer := bufio.NewWriter(create)
	for i := 0; i < len(sort); i++ {
		if sort[i] != "" {
			writer.WriteString(sort[i] + "\n")
		}
	}
	writer.Flush()
	fmt.Println("排序数据重新写入完成", time.Now().Unix())
}

// 二分插入排序
func binInsertSort(arr []string) []string {
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
func findInsertIndex(arr []string, start, end, current int) int {
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
func optimizeQuickSort(arr []string) []string {
	// todo 数据量小于100时插入排序效率更高
	if len(arr) < 100 {
		return binInsertSort(arr)
	}
	privateOptimizeQuickSort(arr, 0, len(arr)-1)
	return arr
}

// 交换
func swap(arr []string, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// 优化快速排序
func privateOptimizeQuickSort(arr []string, left, right int) {
	if right-left < 100 { //截取数组的一段数据 并进行插入排序
		sortPartOfArray(arr, left, right)
	} else {
		// todo 可以极大提升效率
		swap(arr, left, rand.Int()%(right-left+1)+left) //任何一个位置，交换到第一个
		tmp := arr[left]                                // 快速排序 锚点值
		lt := left                                      // 定义使得arr[left+1,lt] 范围内的数据都<tmp lt++
		gt := right + 1                                 // 定义使得arr[gt,right] 范围内的数据都>tmp gt--
		i := left + 1                                   // 定义使得arr[lt+1,i] 范围内的数据都=tmp i++
		for i < gt {
			if arr[i] > tmp {
				swap(arr, i, gt-1)
				gt--
			} else if arr[i] < tmp {
				swap(arr, i, lt+1)
				i++
				lt++
			} else {
				i++
			}
		}
		swap(arr, left, lt)
		if lt-1 > left {
			privateOptimizeQuickSort(arr, left, lt-1)
		}
		privateOptimizeQuickSort(arr, gt, right)
	}
}

// 截取数组的一段数据 并进行插入排序
func sortPartOfArray(arr []string, start, end int) []string {
	if end-start <= 1 {
		// 判断大小 根据排序要求交换位置
		if arr[start] > arr[end] {
			swap(arr, start, end)
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

// 统计密码出现次数
func CountQQPwdTimes(source, save string) {
	open, err := os.Open(source)
	defer open.Close()
	if err != nil {
		fmt.Println("file open error")
		return
	}
	reader := bufio.NewReader(open)
	create, err := os.Create(save)
	defer create.Close()
	if err != nil {
		fmt.Println("file create error")
		return
	}
	writer := bufio.NewWriter(create)
	tmp := ""
	times := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		pwd := string(line)
		if tmp == "" || pwd != tmp {
			if tmp != "" {
				writer.WriteString(tmp + " # " + strconv.Itoa(times) + "\n")
			}
			tmp = pwd
			times = 0
		}
		times++
	}
	writer.Flush()
}

type QQPwdStruct struct {
	pwd   string
	times int
}

// 创建qq密码字典
func CreateQQPasswordDictionary(source, save string) {
	number := Utils.GetFileLineNumber(source)
	fmt.Println(number)
	//os.Exit(1)
	//number := 117334647 //qq
	//number := 4038251     //csdn
	fmt.Println("开始", time.Now().Unix())
	pwd := make([]QQPwdStruct, number, number)
	fmt.Println("数组大小：", len(pwd))
	// 读取文件到结构体
	pwd = readQQPwdToStruct(pwd, source)
	fmt.Println("密码字典写入内存完成", time.Now().Unix())
	fmt.Println("数组大小：", len(pwd))
	// 结构体排序 从多到少
	pwd = qqPwdStructSort(pwd)
	fmt.Println("密码字典排序完成", time.Now().Unix())
	fmt.Println("数组大小：", len(pwd))
	// 排序数据重新写入文件
	qqPwdStructSaveToFile(pwd, save)
	fmt.Println("密码字典写入文件完成", time.Now().Unix())
}

// 读取文件到结构体
func readQQPwdToStruct(pwd []QQPwdStruct, path string) []QQPwdStruct {
	open, _ := os.Open(path)
	defer open.Close()
	reader := bufio.NewReader(open)
	i := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		if tmp != "" {
			split := strings.Split(tmp, " # ")
			if len(split) == 2 {
				pwd[i].pwd = split[0]
				n, _ := strconv.Atoi(split[1])
				pwd[i].times = n
			}
		}
		i++
	}
	return pwd
}

// 结构体排序
func qqPwdStructSort(pwd []QQPwdStruct) []QQPwdStruct {
	if len(pwd) < 100 {
		return qqPwdStructInsertSort(pwd)
	}
	qqPwdStructQuickSort1(pwd, 0, len(pwd)-1)
	return pwd
}

func qqPwdStructQuickSort1(pwd []QQPwdStruct, start, end int) {
	if end-start <= 100 {
		qqPwdStructQuickSort2(pwd, start, end)
	} else {
		swapQQStruct(pwd, start, rand.Int()%(end-start+1)+start)
		tmp := pwd[start] // 中间锚点
		lt := start       // pwd[start+1,lt] >tmp.times lt++
		gt := end + 1     // pwd[gt,end] <tmp.times gt--
		i := start + 1    // pwd[lt+1,i] == tmp.times i++
		for i < gt {
			if pwd[i].times > tmp.times {
				swapQQStruct(pwd, i, lt+1)
				lt++
				i++
			} else if pwd[i].times < tmp.times {
				swapQQStruct(pwd, i, gt-1)
				gt--
			} else {
				i++
			}
		}
		swapQQStruct(pwd, start, lt)
		var wg sync.WaitGroup
		wg.Add(2)
		if lt > start+1 {
			go func() {
				qqPwdStructQuickSort1(pwd, start, lt-1)
				wg.Done()
			}()
		}
		if end > gt {
			go func() {
				qqPwdStructQuickSort1(pwd, gt, end)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

// 插入排序
func qqPwdStructInsertSort(pwd []QQPwdStruct) []QQPwdStruct {
	if len(pwd) <= 1 {
		return pwd
	}
	for i := 1; i < len(pwd); i++ {
		index := findQQPwdStructInsertIndex(pwd, 0, i-1, i)
		if i != index {
			for j := i; j > index; j-- {
				pwd[j], pwd[j-1] = pwd[j-1], pwd[j]
			}
		}
	}
	return pwd
}

// 截取一段数据 并排序
func qqPwdStructQuickSort2(pwd []QQPwdStruct, start, end int) []QQPwdStruct {
	if end-start <= 1 {
		if pwd[start].times < pwd[end].times {
			swapQQStruct(pwd, start, end)
		}
	}
	for i := start + 1; i <= end; i++ {
		index := findQQPwdStructInsertIndex(pwd, start, i-1, i)
		if i != index {
			for j := i; j > index; j-- {
				pwd[j], pwd[j-1] = pwd[j-1], pwd[j]
			}
		}
	}
	return pwd
}

// 交换位置
func swapQQStruct(pwd []QQPwdStruct, i, j int) {
	pwd[i], pwd[j] = pwd[j], pwd[i]
}

// 获取要插入的位置
func findQQPwdStructInsertIndex(pwd []QQPwdStruct, start, end, index int) int {
	if start >= end {
		if pwd[start].times < pwd[index].times {
			return start
		} else {
			return index
		}
	}
	mid := (start + end) / 2
	if pwd[mid].times == pwd[index].times {
		return mid
	} else if pwd[mid].times < pwd[index].times {
		return findQQPwdStructInsertIndex(pwd, start, mid, index)
	} else {
		return findQQPwdStructInsertIndex(pwd, mid+1, end, index)
	}
}

// 排序数据重新写入文件
func qqPwdStructSaveToFile(pwd []QQPwdStruct, save string) {
	create, _ := os.Create(save)
	defer create.Close()
	writer := bufio.NewWriter(create)
	for _, v := range pwd {
		itoa := strconv.Itoa(v.times)
		writer.WriteString(v.pwd + " # " + itoa + "\n")
	}
	writer.Flush()
}

// 获取所有都qq号码
func GetQQNumber(source, save string) {
	open, err := os.Open(source)
	if err != nil {
		fmt.Println("file open err")
		return
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	create, err := os.Create(save)
	if err != nil {
		fmt.Println("file create err")
		return
	}
	defer create.Close()
	writer := bufio.NewWriter(create)
	reg1 := `^[1-9][0-9]{4,11}$`
	regx1 := regexp.MustCompile(reg1)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		split := strings.Split(tmp, "----")
		if len(split) == 2 {
			if regx1.Match([]byte(split[0])) {
				fmt.Fprintln(writer, split[0])
			}
		}
	}
	writer.Flush()
	fmt.Println("所有qq号码提取并写入完成")
}

// 根据qq号码创建索引
// 一次读取每行文件都字节数 然后把数据转换为位置来节省内存
func MakeSearchIndexForQQNumber(source string) []int {
	//number := Utils.GetFileLineNumber(source)
	number := 156783175
	index := make([]int, number, number)
	open, err := os.Open(source)
	Utils.ManageError(err)
	reader := bufio.NewReader(open)
	index[0] = 0
	i := 1
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		tmp := string(line)
		if len(tmp) <= 0 {
			continue
		}
		if i < number {
			// 加上换行符的长度
			index[i] = len(tmp) + 1
		}
		i++
	}
	open.Close()
	fmt.Println("1、数据索引加载完成")
	for j := 0; j < len(index)-1; j++ {
		index[j+1] += index[j]
	}
	fmt.Println("2、数据索引叠加完成")
	return index
}

func GetSearchIndexLineNumber(qq string, source string) {
	indexList := MakeSearchIndexForQQNumber(source)
	open, _ := os.Open(source)
	defer open.Close()
	left := 0
	right := len(indexList) - 1
	for left <= right {
		mid := (left + right) / 2
		index := privateGetDataFromSearchIndexByIndex(indexList, mid, open)
		if index > qq {
			right = mid - 1
		} else if index < qq {
			left = mid + 1
		} else {
			fmt.Println("所在行数是：", mid)
			return
		}
	}
	fmt.Println("位置不存在")
}
func privateGetDataFromSearchIndexByIndex(indexList []int, index int, open *os.File) string {
	open.Seek(0, 0)
	i := indexList[index]
	open.Seek(int64(i), 0)
	bytes := make([]byte, 12, 12)
	read, _ := open.Read(bytes)
	tmp := 0
	for j := 0; j < read; j++ {
		if bytes[j] == '\n' && j >= 5 {
			tmp = j
			break
		}
	}
	return string(bytes[:tmp])
}
