package main

import (
	"fmt"
	"strconv"
	"time"
)

func producer(out chan string) {
	i := 0
	for {
		time.Sleep(time.Second)
		out <- "hi " + strconv.Itoa(i)
		i++
	}
}

func consumer(in chan string, stop chan struct{}) {
	for {
		select {
		case msg := <-in:
			fmt.Println(msg)
		case <-stop:
			fmt.Println("stop")
			return
		}
	}
}

func main() {
	chmsg := make(chan string)
	stop := make(chan struct{})

	go producer(chmsg)
	go consumer(chmsg, stop)
	time.Sleep(5 * time.Second)
	close(stop)
	time.Sleep(5 * time.Second)
}
