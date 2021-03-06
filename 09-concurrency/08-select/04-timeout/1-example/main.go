package main

import (
	"fmt"
	"math/rand"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func server1(ch chan string) {
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	ch <- "hello, it's-a me, server 1"
}

func server2(ch chan string) {
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	ch <- "hello, it's-a me, server 2"
}

func main() {

	rand.Seed(time.Now().UnixNano())
	ch1 := make(chan string)
	ch2 := make(chan string)

	go server1(ch1)
	go server2(ch2)

	select {
	case res := <-ch1:
		fmt.Printf("response: %v \t time: %v \n", res, time.Since(startTime))
		return
	case res := <-ch2:
		fmt.Printf("response: %v \t time: %v \n", res, time.Since(startTime))
		return
	case <-time.After(1000 * time.Millisecond):
		fmt.Printf("timeout \t time: %v \n", time.Since(startTime))
		return
	}

}
