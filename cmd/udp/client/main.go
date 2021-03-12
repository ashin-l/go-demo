package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		i := 0
		for {
			buf := make([]byte, 1024)
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				log.Println("read error: ", err)
			}
			log.Printf("receive msg: %s, from %s\n", buf[:n], addr)
			time.Sleep(3 * time.Second)
			i++
			_, err = conn.WriteTo([]byte("data "+strconv.Itoa(i)), addr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	addr := ""
	fmt.Scanln(&addr)
	log.Println(addr)

	dst, err := net.ResolveUDPAddr("udp", addr+":9093")
	if err != nil {
		log.Fatal(err)
	}

	// The connection can write data to the desired address.
	_, err = conn.WriteTo([]byte("data"), dst)
	if err != nil {
		log.Fatal(err)
	}
	stop := make(chan struct{})
	<-stop
}
