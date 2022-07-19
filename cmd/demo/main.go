package main

import (
	"fmt"
	"os"

	"github.com/ashin-l/go-demo/pkg/db"
	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
)

func main() {
	opt := option.New()
	err := opt.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	logger.Init(opt)
	db.Init(opt)
}
