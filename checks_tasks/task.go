package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	a1 := make([]int, 0, 10)
	a1 = append(a1, []int{1, 2, 3, 4, 5}...)
	fmt.Println(len(a1), cap(a1))
	a2 := append(a1, 6)
	a3 := append(a1, 7)
	fmt.Println("a1: ", (*reflect.SliceHeader)(unsafe.Pointer(&a1)))
	fmt.Println("a2: ", (*reflect.SliceHeader)(unsafe.Pointer(&a2)))
	fmt.Println("a3: ", (*reflect.SliceHeader)(unsafe.Pointer(&a3)))

	fmt.Println(a1, a2, a3)
	fmt.Println(&a1[0], &a2[0], &a3[0])
}
