package Demo

import (
	"sync"
)

// 多线程快速排序
func MultithreadingQuickSortDemo(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	lt, gt, eq := []int{}, []int{}, []int{}
	eq = append(eq, arr[0])
	tmp := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] == tmp {
			eq = append(eq, arr[i])
		} else if arr[i] > tmp {
			gt = append(gt, arr[i])
		} else {
			lt = append(lt, arr[i])
		}
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		lt = MultithreadingQuickSortDemo(lt)
		wg.Done()
	}()
	go func() {
		gt = MultithreadingQuickSortDemo(gt)
		wg.Done()
	}()
	wg.Wait()
	var res []int
	if len(gt) > 0 {
		res = append(res, gt...)
	}
	res = append(res, eq...)
	if len(lt) > 0 {
		res = append(res, lt...)
	}
	return res
}
