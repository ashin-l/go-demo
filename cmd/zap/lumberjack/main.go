package main

import "github.com/ashin-l/go-demo/pkg/zap/lumberjackv2.go"

func main() {
	lumberjackv2.New()
	logger := lumberjackv2.Log
	logger.Info("test")
}
