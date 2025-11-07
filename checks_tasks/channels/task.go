package main

import (
	"fmt"
	"math/rand"
	"time"
)

// имеется функция, которая работает неопределенно долго (до 100 секунд)
func randomTimeWork() int {
	fmt.Println("random time work")
	time.Sleep(time.Duration(rand.Intn(1)) * time.Second)
	fmt.Println("random time work done")
	return 100
}

// написать обертку для этой функции, которая будет прерывать выполнение, если
// функция работает дольше 3 секунд, и возвращать ошибку
func predictableTimeWork(s int) {
	ch := make(chan int)

	go func() {
		r := randomTimeWork()
		ch <- r
		close(ch)
	}()

	select {
	case r := <-ch:
		fmt.Println("predictableTimeWork completed", r)
	case t := <-time.After(time.Duration(s) * time.Second):
		fmt.Println("timed out", t)
	}
}

func main() {
	fmt.Println("main")
	predictableTimeWork(3)
	fmt.Println("main completed")
}
