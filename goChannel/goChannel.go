package main

import (
	"fmt"
	"time"
)

func main() {
	// The buffer size is the number of elements that can be sent to the channel without the send blocking.
	// By default, a channel has a buffer size of 0 (you get this with make(chan int)).
	// This means that every single send will block until another goroutine receives from the channel.
	// A channel of buffer size 1 can hold 1 element until sending blocks,
	// in this case channelSize = 10
	buf := make(chan int, 10)

	go func() {
		for i := 0; i < 50; i++ {
			fmt.Println("<-- ", i)
			buf <- i
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			msg := <-buf
			fmt.Println("--> ", msg)
			time.Sleep(time.Second * 5)
		}
	}()

	var input string
	fmt.Scanln(&input)
}

// hamed@siasi.co.uk
