package main

import (
	"fmt"
	"os/signal"
	"strconv"
	"time"

	//import the Paho Go MQTT library
	"os"

	"github.com/ashin-l/go-demo/pkg/util"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/goccy/go-json"
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
	defQos      = "1"

	envAddr     = "MY_ADDR"
	envTopic    = "MY_TOPIC"
	envUsername = "MY_USERNAME"
	envPassword = "MY_PASSWORD"
	envClientid = "MY_CLIENTID"
	envQos      = "MY_QOS"
)

type config struct {
	addr     string
	topic    string
	username string
	password string
	clientid string
	qos      int
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s, MSGID: %d\n", msg.Topic(), msg.MessageID())
	data := make(map[string]interface{})
	json.Unmarshal(msg.Payload(), &data)

	// ts, _ := strconv.ParseInt(data["ts"].(int64), 0, 64)
	ts := int64(data["ts"].(float64))
	fmt.Println(time.UnixMilli(ts))
}

func loadConfig() config {
	qos, err := strconv.Atoi(util.Env(envQos, defQos))
	if err != nil {
		fmt.Println("error", envQos)
		os.Exit(0)
	}
	return config{
		addr:     util.Env(envAddr, defAddr),
		topic:    util.Env(envTopic, defTopic),
		username: util.Env(envUsername, defUsername),
		password: util.Env(envPassword, defPassword),
		clientid: util.Env(envClientid, defClientid),
		qos:      qos,
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
	fmt.Println(cfg.username)
	fmt.Println(cfg.topic)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	if token := c.Subscribe(cfg.topic, byte(cfg.qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	c.Disconnect(250)
}
