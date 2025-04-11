package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ashin-l/go-demo/pkg/kafka"
	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
)

const (
	msg   = "go-demo test msg"
)

func HandleMsg(in chan []byte, stop chan struct{}) {
	go func() {
		logger.Logger().Info("开始接收kafka消息")
		for {
			select {
			case v := <-in:
				logger.Logger().Info("接收到kafka消息:", string(v))
			case <-stop:
				logger.Logger().Info("停止接收周界数据")
				return
			}
		}
	}()

}

func WaitSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	s := <-signals
	logger.Logger().Info("signal:", s.String())
}

func main() {
	opt := option.New()
	err := opt.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Init(opt)
	err = kafka.Init(opt)
	if err != nil {
		logger.Logger().Error(err)
		os.Exit(1)
	}
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	kmc := make(chan []byte)
	ctx := context.Background()
	kafka.Sub(ctx, map[string]chan []byte{opt.Kafka.Topic: kmc})
	HandleMsg(kmc, stop)
	pmc, err := kafka.Pub(wg, opt.Kafka.Topic)
	if err != nil {
		logger.Logger().Error(err)
		os.Exit(1)
	}
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				pmc <- []byte(msg)
			case <-stop:
				close(pmc)
				return
			}
		}
	}()
	WaitSignal()
	close(stop)
	ctx.Done()
	wg.Wait()
	kafka.Stop()
	logger.Logger().Info("----------- down ---------------")
}