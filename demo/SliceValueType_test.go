package demo

import (
	"fmt"
	"testing"
)

/**
问题：
1、通过 append() 函数追加不会影响原来的值，这个我能理解，因为是值传递。
2、对函数进行值传递，修改了下标后，为什么会对原来的 slice 造成影响？
3、在函数内打印 a 变量的内存地址还是原来的内存地址？我理解的是值传递是对数据的拷贝，应该指向另外的内存地址吧？
*/

func TestDeref(t *testing.T) {
	var si []int
	si = append(si, 1, 2, 3, 4, 5)
	fmt.Println(len(si), cap(si)) // 5 6
	fmt.Printf("%p \n", si)       // 0xc0000ac060
	Deref(si)
	fmt.Println(si)               // [1 10 3 4 5]
	fmt.Printf("%p \n", si)       // 0xc0000ac060
	fmt.Printf("%p \n", &si)      // 0xc00000e2a0
	fmt.Println(len(si), cap(si)) // 5 6
}

// 参数 a 是：值传递
func Deref(a []int) {
	// ############## 值传递后， si 变量内存地址还是原来的地址
	fmt.Printf("%p \n", a)      // 0xc0000ac060
	fmt.Println(a)              // [1 2 3 4 5]
	a[1] = 10                   // ####### 这里修改了下标 1 却会影响外面的值，不是说值传递是传递的是副本吗？
	a = append(a, 6)            // ##### 这里追加不会影响原来的值
	fmt.Println(a)              // [1 10 3 4 5 6]
	fmt.Printf("%p \n", a)      // 0xc0000ac060
	fmt.Printf("%p \n", &a)     // 0xc00000e270
	fmt.Println(len(a), cap(a)) //	6 6
}
