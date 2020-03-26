package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	Addr     string
	PassWord string
	DB       int
}

var c *redis.Client

func New(m Redis) (err error) {
	c = redis.NewClient(&redis.Options{
		Addr:     m.Addr,
		Password: m.PassWord,
		DB:       m.DB,
	})

	pong, err := c.Ping().Result()
	fmt.Println(pong, err)
	return
}
