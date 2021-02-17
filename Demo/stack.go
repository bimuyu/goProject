package Demo

import (
	"container/list"
	"fmt"
)

// 简单栈实现
func SimpleStackDemo() {
	l := list.New()
	for i := 0; i < 9; i++ {
		l.PushBack(i)
	}
	fmt.Println("长度为：", l.Len())
	for l.Len() != 0 {
		back := l.Back()
		l.Remove(back)
		fmt.Println(back.Value.(int))
	}
}

// 简单归并实现
func SimpleArrayMerge(arr []string) {
	l := list.New()

	for i := 0; i < len(arr); i++ {
		l.PushBack(arr[i])
	}

	for l.Len() != 1 {
		n1 := l.Back()
		l.Remove(n1)

		n2 := l.Back()
		l.Remove(n2)

		s1 := n1.Value.(string)
		s2 := n2.Value.(string)

		n := s1 + s2
		l.PushBack(n)
	}

	v := l.Back()
	l.Remove(v)
	s := v.Value.(string)
	fmt.Println(s)
}
