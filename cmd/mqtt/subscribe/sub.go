package main

import (
	"fmt"
	"os/signal"

	//import the Paho Go MQTT library
	"os"

	"github.com/ashin-l/go-demo/pkg/util"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"strconv"
	//"github.com/astaxie/beego/config"
)

const (
	defAddr = "tcp://192.168.152.44:1883"
	//defTopic    = "rbtest"
	defTopic    = "v1/pm/response"
	defUsername = ""
	defPassword = ""
	defClientid = "test-cli-1"

	envAddr     = "MY_ADDR"
	envTopic    = "MY_TOPIC"
	envUsername = "MY_USERNAME"
	envPassword = "MY_PASSWORD"
	envClientid = "MY_CLIENTID"
)

type config struct {
	addr     string
	topic    string
	username string
	password string
	clientid string
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func loadConfig() config {
	return config{
		addr:     util.Env(envAddr, defAddr),
		topic:    util.Env(envTopic, defTopic),
		username: util.Env(envUsername, defUsername),
		password: util.Env(envPassword, defPassword),
		clientid: util.Env(envClientid, defClientid),
	}
}

func main() {
	cfg := loadConfig()

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(cfg.addr)
	opts.SetClientID(cfg.clientid)
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername(cfg.username)
	opts.SetPassword(cfg.password)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	if token := c.Subscribe(cfg.topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
	<-sigChan
	c.Disconnect(250)
}
