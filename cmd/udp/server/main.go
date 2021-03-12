package main

import (
	"log"
	"net"
)

func main() {
	log.Println("server start...")
	l, err := net.ListenPacket("udp", "0.0.0.0:9090")
	if err != nil {
		log.Println("server error: ", err.Error())
	}
	defer l.Close()
	for {
		buf := make([]byte, 1024)
		n, addr, err := l.ReadFrom(buf)
		if err != nil {
			log.Println("read error: ", err)
		}
		log.Printf("receive msg: %s, from %s\n", buf[:n], addr)
		_, err = l.WriteTo([]byte("receive"), addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
