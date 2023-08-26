package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func RunOnce() {
	var once sync.Once

	onceBody := func() {
		fmt.Println("我只执行一次")
	}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// OnceFunc返回一个可以并发调用的函数，它可以被调用多次。
// 即使返回的函数被调用多次，f也只会被调用一次。
func RunOnceFunc() {
	onceBody := func() {
		fmt.Println("我只执行一次")
	}

	foo := sync.OnceFunc(onceBody)
	for i := 0; i < 10; i++ {
		foo()
	}
}

// OnceValue返回一个函数， 这个函数会返回f的返回值。多次调用都会返回同一个值。
func GetOnceValue() {
	randvalue := func() int {
		return rand.Int()
	}

	bar := sync.OnceValue(randvalue)
	for i := 0; i < 10; i++ {
		fmt.Println("打印随机值:", randvalue())
		fmt.Println("多次调用都会返回同一个值: ", bar())
	}
}
