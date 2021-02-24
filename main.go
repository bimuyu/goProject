package main

import (
	"bigData/Demo"
	"bigData/QQ"
	"bigData/Utils"
	"fmt"
	"os"
	"time"
)

func main() {
	/* ------------ qq密码字典开始 ------------- */
	// 一、提取所有密码
	//GetQQPasswordFromFile()
	// 二、密码排序
	//QQPasswordSort()
	// 三、计算密码次数
	//QQPasswordTimes()
	// 四、按照密码次数排序
	//QQPasswordDictionary()
	// 五、获取所有都qq号
	//getAllQQNumber()
	// 六、qq号码排序
	//QQNumberSort()
	/* ------------ qq密码字典结束 ------------- */
	//mysqlConnectDemo()
	//dataClearToNewFileDemo()
	//binSearchDemo()
	//quickSortDemo()
	//readFileToArrayAndSortDemo()
	//dataUniqueAndCountDemo()
	//dataSortByTimesDemo()
	//insertSortDemo()
	//binInsertSortDemo()
	//optimizeQuickSortDemo()
	//traverseDemo()
	//numberEqualSplit()
	//sliceFile()
	//simpleStackDemo()
	//simpleArrayMergeDemo()
	//arrayMerger()
	//simpleMultithreadingQuickSortDemo()
	//cmdDemo()
	//intAndBytesExchange()
	//floatAndByteExchange()
	//sortInterfaceDemo()
	//Utils.JsonDemo()
	//fileMergeDemo()
	//selectSortDemo()
	//heapSortDemo()
	//bubbleSortDemo()
	//shellSortDemo()
	//radixSortDemo()
	//getQQNumberSearchIndex()
	binSearchQQNumberSearchIndex()
}

func mysqlConnectDemo() {
	host := "127.0.0.1"
	port := "3306"
	password := "123456"
	user := "root"
	dbName := "go_project"
	charset := "utf8"
	connect := Utils.MysqlConnect(host, port, password, user, dbName, charset)
	fmt.Println(connect)
}

func dataClearToNewFileDemo() {
	sourcePath := "/Users/magic/web/golang/source/demo.txt"
	savePath := "/Users/magic/web/golang/source/demo_innid.txt"
	Utils.DataClearToNewFile(sourcePath, savePath)
}

func binSearchDemo() {
	arr := make([]int, 1024*1024, 1024*1024)
	for i := 0; i < 1024*1024; i++ {
		arr[i] = 2*i + 1
	}
	fmt.Println(Utils.BinSearch(arr, 13124))
}

func quickSortDemo() {
	arr := Utils.MakeArrFunc(20)
	fmt.Println(arr)
	fmt.Println(Utils.QuickSort(arr))
}

func readFileToArrayAndSortDemo() {
	var arr []int
	fmt.Println("开始")
	sourcePath := "/Users/magic/web/golang/source/demo_innid.txt"
	sortArr := Utils.ReadFileToArrayAndSortDemo(sourcePath, arr)
	fmt.Println("数据排序完成")
	savePath := "/Users/magic/web/golang/source/demo_innid_sort.txt"
	Utils.SaveArrayToFile(savePath, sortArr)
	fmt.Println("排序数据写入完成")
}

func dataUniqueAndCountDemo() {
	sourcePath := "/Users/magic/web/golang/source/demo_innid_sort.txt"
	savePath := "/Users/magic/web/golang/source/demo_innid_sort_unique.txt"
	Utils.DataUniqueAndCount(sourcePath, savePath)
}

func dataSortByTimesDemo() {
	sourcePath := "/Users/magic/web/golang/source/demo_innid_sort_unique.txt"
	demoStructs := Utils.DataSortByTimes(sourcePath)
	fmt.Println("数据读取到结构体完成")
	sortData := Utils.DemoStructInfoQuickSort(demoStructs)
	fmt.Println("结构体排序完成")
	savePath := "/Users/magic/web/golang/source/demo_innid_sort_times.txt"
	Utils.SaveDemoStructToFile(sortData, savePath)
	fmt.Println("结构体排序数据写入完成")
}

func insertSortDemo() {
	//Utils.InsertSortDemo()
	arr := Utils.MakeArrFunc(20)
	fmt.Println(arr)
	fmt.Println(Utils.SimpleInsertSort(arr))
}

func binInsertSortDemo() {
	//Utils.InsertSortDemo()
	arr := Utils.MakeArrFunc(20)
	//arr := []int{556,7,10,678,155,370,307, 152, 825, 224, 844, 891, 938, 734, 430, 972, 370, 640, 846, 820}
	fmt.Println(arr)
	fmt.Println(Utils.BinInsertSort(arr))
}

func optimizeQuickSortDemo() {
	arr := Utils.MakeArrFunc(30)
	//arr := []int{836, 1, 10, 706, 656, 386, 54, 531, 872, 312, 837, 538, 532, 853, 346, 350, 320, 915, 959, 853, 671, 134, 583, 170, 436, 335, 748, 271, 451, 700}
	//arr := []int{371,771,617,444,187,313,102,997,29,172,717,375,3,447,412,816,237,908,98,991,211,318,502,415,485,266,758,786,832,607}
	//arr := []int{371,771,617,444,187,313,102,997,29,172,717,375,3,447,412,816,237,908,98,991,211,318,502,415,485,266,758,786,832,607}
	fmt.Println(arr)
	Utils.OptimizeQuickSort(arr)
	fmt.Println(arr)
}

func traverseDemo() {
	path := "/Users/magic/web/golang/source"
	Utils.TraverseDir(path)
}

func GetQQPasswordFromFile() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/QQBig.txt"
	savePath := "/Users/magic/web/golang/source/QQPWD/only_qq_pwd2.txt"
	QQ.GetQQPwd(sourcePath, savePath)
}

func QQPasswordSort() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/only_qq_pwd.txt"
	savePath := "/Users/magic/web/golang/source/QQPWD/qq_pwd_sort.txt"
	QQ.QQPwdQuickSort(sourcePath, savePath)
}

func QQPasswordTimes() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/qq_pwd_sort.txt"
	savePath := "/Users/magic/web/golang/source/QQPWD/qq_pwd_times.txt"
	//sourcePath := "/Users/magic/web/golang/source/CSDNPWD/csdn_pwd_sort.txt"
	//savePath := "/Users/magic/web/golang/source/CSDNPWD/csdn_pwd_times.txt"
	QQ.CountQQPwdTimes(sourcePath, savePath)
}

func QQPasswordDictionary() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/qq_pwd_times.txt"
	savePath := "/Users/magic/web/golang/source/QQPWD/qq_pwd_dictionary.txt"
	//sourcePath := "/Users/magic/web/golang/source/CSDNPWD/csdn_pwd_times.txt"
	//savePath := "/Users/magic/web/golang/source/CSDNPWD/csdn_pwd_dictionary.txt"
	QQ.CreateQQPasswordDictionary(sourcePath, savePath)
}

func numberEqualSplit() {
	data := Utils.EqualSplitData(97, 10)
	fmt.Println(data, 97)
	data = Utils.EqualSplitData(101, 10)
	fmt.Println(data, 101)
	data = Utils.EqualSplitData(100, 10)
	fmt.Println(data, 102)
}

func sliceFile() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/QQBig.txt"
	savePath := "/Users/magic/web/golang/source/split_file/QQBig_"
	Utils.SliceFileToSmall(sourcePath, savePath, 9)
}

func simpleStackDemo() {
	Demo.SimpleStackDemo()
}

func simpleArrayMergeDemo() {
	arr := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	Demo.SimpleArrayMerge(arr)
}

func arrayMerger() {
	arr := Utils.MakeArrFunc(20)
	fmt.Println(arr)
	sort := Utils.ArrayMergeSort(arr)
	fmt.Println(sort)
}

func simpleMultithreadingQuickSortDemo() {
	arr := Utils.MakeArrFunc(20)
	fmt.Println(arr)
	demo := Demo.MultithreadingQuickSortDemo(arr)
	fmt.Println(demo)
}

func cmdDemo() {
	cmd := "ls"
	Demo.CmdDemo(cmd)
}

func intAndBytesExchange() {
	n := 99
	bytes := Utils.IntToBytes(n)
	toInt := Utils.BytesToInt(bytes)
	fmt.Println(n)
	fmt.Println(bytes)
	fmt.Println(toInt)
}

func floatAndByteExchange() {
	n := 13.424241223
	fmt.Println(float32(n))
	fmt.Println(Utils.BytesToFloat32(Utils.Float32ToBytes(float32(n))))
	fmt.Println(Utils.BytesToFloat64(Utils.Float64ToBytes(n)))
}

func sortInterfaceDemo() {
	data := new(Utils.QuickSortStruct)
	/*n1 := Utils.SortDemoStruct{"abc", 0}
	n2 := Utils.SortDemoStruct{"abc", 1}
	n3 := Utils.SortDemoStruct{"abc", 2}
	n4 := Utils.SortDemoStruct{"abc", 3}
	n5 := Utils.SortDemoStruct{"abc", 4}
	n6 := Utils.SortDemoStruct{"abc", 5}
	n7 := Utils.SortDemoStruct{"abc", 6}
	n8 := Utils.SortDemoStruct{"abc", 7}
	n9 := Utils.SortDemoStruct{"abc", 8}
	n10 := Utils.SortDemoStruct{"abc", 9}
	n11 := Utils.SortDemoStruct{"abc", 10}
	n12 := Utils.SortDemoStruct{"abc", 11}
	n13 := Utils.SortDemoStruct{"abc", 12}
	data.Data = []interface{}{n1, n2, n3, n4, n5, n6, n7, n8, n9, n10, n11, n13, n12}*/
	//data.Data = []interface{}{7, 9, 2, 8, 3, 3, 3, 9, 9, 11, 17, 16, 13}
	data.Data = []interface{}{"a", "b", "c", "d", "e", "f", "abc", "abb"}
	data.IsAsc = true
	data.DataType = "string"
	data.IsFile = true
	data.InitMyFunc()

	fmt.Println("排序前：", data.Data)
	data.Sort()
	fmt.Println("排序后：", data.Data)
}

func fileMergeDemo() {
	fileList := []string{
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_0.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_1.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_2.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_3.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_4.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_5.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_6.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_7.txt",
		"/Users/magic/web/golang/source/csdn_split/csdn_pwd_8.txt",
	}
	path := Utils.MergeFileListAsOne(fileList)
	fmt.Println(path)
}

func selectSortDemo() {
	arrFunc := Utils.MakeArrFunc(10)
	fmt.Println(arrFunc)
	sort := Utils.SimpleSelectSort(arrFunc, true)
	fmt.Println(sort)
}

func heapSortDemo() {
	arrFunc := Utils.MakeArrFunc(10)
	fmt.Println(arrFunc)
	sort := Utils.SimpleHeapSort(arrFunc)
	fmt.Println(sort)
}

func bubbleSortDemo() {
	arrFunc := Utils.MakeArrFunc(10)
	fmt.Println(arrFunc)
	sort := Utils.BubbleSortDemo(arrFunc, false)
	fmt.Println(sort)
}

func shellSortDemo() {
	arrFunc := Utils.MakeArrFunc(30)
	fmt.Println(time.Now().Unix())
	Utils.ShellSortDemo(arrFunc)
	fmt.Println(time.Now().Unix())
	Utils.OptimizeShellSort(arrFunc)
	fmt.Println(time.Now().Unix())
}

func radixSortDemo() {
	arr := []int{1, 2, 3, 4, 4, 3, 2, 1, 2, 3, 1, 2, 3, 4, 1, 2, 2, 2, 2, 3, 4, 1}
	fmt.Println(Utils.RadixSortDemo(arr))
}

func getAllQQNumber() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/QQBig.txt"
	savePath := "/Users/magic/web/golang/source/QQPWD/only_qq_number.txt"
	QQ.GetQQNumber(sourcePath, savePath)
}

func QQNumberSort() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/only_qq_number.txt"
	savePath := "/Users/magic/web/golang/source/QQPWD/qq_number_sort.txt"
	QQ.QQPwdQuickSort(sourcePath, savePath)
}

func getQQNumberSearchIndex() {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/qq_number_sort.txt"
	searchIndex := QQ.MakeSearchIndexForQQNumber(sourcePath)
	open, _ := os.Open(sourcePath)
	// 移动位置
	open.Seek(0, 0)
	for {
		var lineNumber int64
		fmt.Scanf("%d", &lineNumber)
		fmt.Println("位置：",searchIndex[lineNumber])
		open.Seek(int64(searchIndex[lineNumber]), 0)
		bytes := make([]byte, 12, 12)
		read, _ := open.Read(bytes)
		var endPosition int
		for i := 0; i < read; i++ {
			// 换行符表示结束 5是最短的qq号
			if bytes[i] == '\n' && i >= 5 {
				endPosition = i
				break
			}
		}
		fmt.Println("对应的数据：", string(bytes[:endPosition]))
	}
	open.Close()
}

func binSearchQQNumberSearchIndex()  {
	sourcePath := "/Users/magic/web/golang/source/QQPWD/qq_number_sort.txt"
	qq := "9324834"
	QQ.GetSearchIndexLineNumber(qq, sourcePath)
}
