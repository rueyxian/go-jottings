package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {

	fmt.Println("[before] num of goroutines: ", runtime.NumGoroutine())
	fmt.Println("begin operation...")
	fmt.Println()
	operation()
	time.Sleep(time.Second)
	fmt.Println()
	fmt.Println("end operation...")
	fmt.Println("[after] num of goroutines: ", runtime.NumGoroutine())

}

func operation() {
	done := make(chan struct{}, 2)
	ch1 := pipe("A", "B", "C")
	ch2 := pipe("1", "2", "3")

	output := merge(done, ch1, ch2)
	// ==============================

	// for o := range output {
	//   fmt.Println("receive:", o)
	// }

	// to simulate missing some <-output
	loop := 5
	for i := 0; i < loop; i++ {
		fmt.Println("        recv:", <-output)
	}

	// if there is a blocking – channel's receiver is waiting for sender,
	// send "done" signal to unblock and move on.
	// again, here's the problem, we don't know how much "done" signal do we need.
	// Here we are assuming there are potentially one blocking on each channel
	done <- struct{}{}
	done <- struct{}{}

	return

}

// ============================================================
func pipe(strs ...string) <-chan string {
	ret := make(chan string, len(strs))
	for _, str := range strs {
		ret <- str
	}
	close(ret)
	return ret
}

// ============================================================
func merge(done <-chan struct{}, chs ...<-chan string) <-chan string {

	ret := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(chs))

	output := func(ch <-chan string) {
		defer wg.Done()
		for s := range ch {
			select {
			case ret <- s:
				fmt.Println("send:", s)
			case <-done:
			}
		}
	}

	for _, ch := range chs {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(ret)
	}()

	return ret
}
