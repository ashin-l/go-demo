package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ashin-l/go-demo/pkg/common"
	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/mqtt"
	"github.com/ashin-l/go-demo/pkg/option"
)

func main() {
	opt := option.New()
	err := opt.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Init(opt)
	err = mqtt.Init(opt)
	if err != nil {
		logger.Logger().Error(err)
		os.Exit(1)
	}

	logger.Logger().Infof("opt: %v", opt)
	wg := &sync.WaitGroup{}
	stop := make(chan struct{})
	for _, val := range opt.MqttTopics {
		mc := make(chan []byte, 1)
		mqtt.Pub(wg, val.Topic, val.Qos, mc)
		go sendMessage(val, mc, stop)
	}
	common.WaitSignal()
	close(stop)
	wg.Wait()
	logger.Logger().Info("----------- down ---------------")
}

func sendMessage(mt option.MqttTopic, mc chan []byte, stop chan struct{}) {
	ticker := time.NewTicker(time.Duration(mt.Interval) * time.Second)
	ts := time.Now().Unix()
	i := 1
	for {
		select {
		case <-ticker.C:
			ts += int64(mt.Interval)
			payload := fmt.Sprintf(mt.Fmtstr, ts)
			mc <- []byte(payload)
			logger.Logger().Infow("publish", "topic", mt.Topic, "index", i)
			i++
		case <-stop:
			logger.Logger().Info("stop publish ", mt.Topic)
			close(mc)
			return
		}
	}
}
