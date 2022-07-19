package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	timeout := time.Duration(3 * time.Second)
	_, err := net.DialTimeout("tcp", "192.168.61.11:554", timeout)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	fmt.Println("tcp server is ok")
}
