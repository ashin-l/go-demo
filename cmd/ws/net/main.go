package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

var reply string

func upper(ws *websocket.Conn) {
	var err error
	if err = websocket.Message.Send(ws, strings.ToUpper(reply)); err != nil {
		fmt.Println(err)
	}
	for {

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println(err)
			continue
		}

		if err = websocket.Message.Send(ws, strings.ToUpper(reply)); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func handleHi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	str := "himan"
	w.Write([]byte(str))
}

func main() {
	if len(os.Args) != 4 {
		os.Exit(0)
	}
	path := os.Args[1]
	port := os.Args[2]
	reply = os.Args[3]
	http.HandleFunc(path+"/hi", handleHi)
	http.Handle(path+"/ws", websocket.Handler(upper))

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
