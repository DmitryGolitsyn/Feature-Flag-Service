package main

import (
	"fmt"
	"time"
)

func reader(chIn <-chan int) {
	go func() {
		fmt.Println("reader started")
		for v := range chIn {
			time.Sleep(500 * time.Millisecond)
			fmt.Println(v)
		}
	}()
}

func double(chIn <-chan int) <-chan int {
	chOut := make(chan int)

	go func() {
		fmt.Println("double started")
		for {
			value, ok := <-chIn
			if !ok {
				break
			}
			chOut <- value * 2
		}
		close(chOut)
	}()

	return chOut
}

func writer() <-chan int {
	chOut := make(chan int)

	go func() {
		fmt.Println("writer started")
		for i := range 10 {
			chOut <- i + 1
		}
		close(chOut)
	}()

	return chOut
}

// Написать 3 функции:
// writer - генерит числа от 1 до 10
// doubler - умножает числа на 2, имитируя работу (500ms)
// reader - читает и выводит на экран
func main() {
	fmt.Println("main")
	reader(double(writer()))
	time.Sleep(10 * time.Second)
}
