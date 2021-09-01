package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://admin:ai123456@192.168.165.125/onvif-http/snapshot?")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("test.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
