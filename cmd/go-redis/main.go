package main

import (
	"fmt"

	"github.com/ashin-l/go-demo/pkg/redis"
)

func main() {
	m := redis.Redis{
		Addr:     "192.168.152.41:6379",
		PassWord: "",
		DB:       0,
	}
	err := redis.New(m)
	if err != nil {
		fmt.Println("new redis client error:", err)
	}
}
