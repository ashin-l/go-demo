package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ashin-l/go-demo/pkg/util"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	defAddr     = "tcp://127.0.0.1:1883"
	defTopic    = "mytest"
	defUsername = "xxx"
	defPassword = "xxx"
	defClientid = ""
	defInterval = "3000"
	defDjson    = `{"ts":%d,"val":"hello world!"}`

	envAddr     = "MY_ADDR"
	envTopic    = "MY_TOPIC"
	envUsername = "MY_USERNAME"
	envPassword = "MY_PASSWORD"
	envClientid = "MY_CLIENTID"
	envInterval = "MY_INTERVAL"
	envDjson    = "MY_DJSON"
)

type config struct {
	addr     string
	topic    string
	username string
	password string
	clientid string
	interval int
	djson    string
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func loadConfig() config {
	iv, err := strconv.Atoi(util.Env(envInterval, defInterval))
	if err != nil {
		fmt.Println("error", envInterval)
		os.Exit(0)
	}
	return config{
		addr:     util.Env(envAddr, defAddr),
		topic:    util.Env(envTopic, defTopic),
		username: util.Env(envUsername, defUsername),
		password: util.Env(envPassword, defPassword),
		clientid: util.Env(envClientid, defClientid),
		interval: iv,
		djson:    util.Env(envDjson, defDjson),
	}
}

func main() {
	cfg := loadConfig()
	fmt.Printf("%v\n", cfg)

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(cfg.addr)
	if cfg.clientid != "" {
		opts.SetClientID(cfg.clientid)
		fmt.Println("clientId: ", cfg.clientid)
	}
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername(cfg.username)
	opts.SetPassword(cfg.password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(false)

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

	if cfg.interval == 0 {
		token := c.Publish(cfg.topic, 0, false, cfg.djson)
		token.Wait()
		time.Sleep(10 * time.Second)
		c.Disconnect(50)
	} else {
		exitchan := make(chan struct{})
		ticker := time.NewTicker(time.Duration(cfg.interval) * time.Millisecond)
		go func() {
			sigChan := make(chan os.Signal)
			signal.Notify(sigChan)
			<-sigChan
			ticker.Stop()
			close(exitchan)
		}()
		mtime := time.Now().UnixNano() / 1e6
		//mtime := time.Now().Unix()
		for {
			select {
			case <-ticker.C:
				mtime += int64(cfg.interval)
				payload := fmt.Sprintf(cfg.djson, mtime)
				token := c.Publish(cfg.topic, 0, false, payload)
				token.Wait()
				fmt.Println("publish")
			case <-exitchan:
				c.Disconnect(30)
				fmt.Println("disconnect!")
				return
			}
		}

	}
}
