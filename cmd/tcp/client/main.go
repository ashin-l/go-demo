package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.165.43:20001")
	if err != nil {
		fmt.Println("Connect to TCP server failed ,err:", err)
		return
	}
	fmt.Println("111111111")
	time.Sleep(10 * time.Second)
	fmt.Println("222222222")
	conn.Close()
}
