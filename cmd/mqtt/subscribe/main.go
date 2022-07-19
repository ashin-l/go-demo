package main

import (
	"fmt"
	"os"

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

	for _, val := range opt.MqttTopics {
		go sub(val)
		logger.Logger().Infow("sub mqtt topic", "topic", val.Topic, "Qos", val.Qos)
	}
	common.WaitSignal()
	logger.Logger().Info("----------- down ---------------")
}

func sub(mt option.MqttTopic) {
	mc, err := mqtt.Sub(mt.Topic, mt.Qos)
	if err != nil {
		logger.Logger().Error(err)
		return
	}
	for msg := range mc {
		logger.Logger().Infow("receive message", "topic", mt.Topic, "msg", string(msg))
	}
}
