package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ashin-l/go-demo/pkg/util"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	sugar *zap.SugaredLogger
	sum   int64
	muSum sync.Mutex
)

const (
	defAddr         = "tcp://127.0.0.1:1883"
	defTopic        = "mytest"
	defUsername     = "xxx"
	defPassword     = "xxx"
	defClientid     = ""
	defInterval     = "3000"
	defUUID         = "xxx"
	defLogLevel     = "info"
	defLogInfoPath  = "logs/info.log"
	defLogErrorPath = "logs/error.log"

	envAddr         = "MY_ADDR"
	envTopic        = "MY_TOPIC"
	envUsername     = "MY_USERNAME"
	envPassword     = "MY_PASSWORD"
	envClientid     = "MY_CLIENTID"
	envInterval     = "MY_INTERVAL"
	envUUID         = "MY_UUID"
	envLogLevel     = "MY_LOGLEVEL"
	envLogInfoPath  = "MY_LOGINFOPATH"
	envLogErrorPath = "MY_LOGERRORPATH"
)

type config struct {
	addr         string
	topic        string
	username     string
	password     string
	clientid     string
	interval     int
	uuid         string
	logLevel     string
	logInfoPath  string
	logErrorPath string
	// djson    string
}

type TlmMeteorological struct {
	Ts            int64   `json:"ts"`
	Uuid          string  `json:"uuid"`
	Temperature   float32 `json:"temperature"`
	Humidity      float32 `json:"humidity"`
	WindDirection float32 `json:"windDirection"`
	WindSpeed     float32 `json:"windSpeed"`
	Pressure      float32 `json:"pressure"`
	Rainfall      float32 `json:"rainfall"`
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	sugar.Info(string(msg.Payload()))
}

func loadConfig() config {
	iv, err := strconv.Atoi(util.Env(envInterval, defInterval))
	if err != nil {
		fmt.Println("error", envInterval)
		os.Exit(0)
	}
	return config{
		addr:         util.Env(envAddr, defAddr),
		topic:        util.Env(envTopic, defTopic),
		username:     util.Env(envUsername, defUsername),
		password:     util.Env(envPassword, defPassword),
		clientid:     util.Env(envClientid, defClientid),
		interval:     iv,
		uuid:         util.Env(envUUID, defUUID),
		logLevel:     util.Env(envLogLevel, defLogLevel),
		logInfoPath:  util.Env(envLogInfoPath, defLogInfoPath),
		logErrorPath: util.Env(envLogErrorPath, defLogErrorPath),
	}
}

func sendMsg(c MQTT.Client, topic string, data *TlmMeteorological) {
	data.Ts = time.Now().UnixNano() / 1e6
	data.Temperature = rand.Float32() + float32(rand.Intn(70))
	data.Humidity = rand.Float32() + float32(rand.Intn(70))
	payload, err := json.Marshal(&data)
	if err != nil {
		sugar.Warn("marshal error:", err)
		return
	}
	token := c.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return
	}
	muSum.Lock()
	sum++
	sugar.Infow("send msg success total", "total", sum)
	muSum.Unlock()
}

func main() {
	cfg := loadConfig()
	fmt.Printf("%v\n", cfg)
	sugar = log.NewLog(cfg.logInfoPath, cfg.logErrorPath, cfg.logLevel).Sugar()
	defer sugar.Sync()
	// sugar.Infow("failed to fetch URL",
	// 	// Structured context as loosely typed key-value pairs.
	// 	"url", "http",
	// 	"attempt", 3,
	// 	"backoff", time.Second,
	// )
	// sugar.Infof("Failed to fetch URL: %s", "http1")
	// sugar.Warnf("warn test")
	// sugar.Error("error test")
	// sugar.Debugf("debug %s", "debug")

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(cfg.addr)
	if cfg.clientid != "" {
		opts.SetClientID(cfg.clientid)
	}
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername(cfg.username)
	opts.SetPassword(cfg.password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(false)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		sugar.Fatal(token.Error())
	}

	// if token := c.Subscribe(cfg.topic, 0, nil); token.Wait() && token.Error() != nil {
	// 	sugar.Fatal(token.Error())
	// }

	exitchan := make(chan struct{})
	data := &TlmMeteorological{
		Uuid:          cfg.uuid,
		WindDirection: 174.32912,
		WindSpeed:     2.3287432,
		Pressure:      1011.3,
		Rainfall:      0,
	}
	ticker := time.NewTicker(time.Duration(cfg.interval) * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				// data.Ts = time.Now().UnixNano() / 1e6
				// data.Temperature = rand.Float32() + float32(rand.Intn(70))
				// data.Humidity = rand.Float32() + float32(rand.Intn(70))
				// payload, err := json.Marshal(&data)
				// if err != nil {
				// 	sugar.Warn("marshal error:", err)
				// 	continue
				// }
				// token := c.Publish(cfg.topic, 0, false, payload)
				// token.Wait()
				go sendMsg(c, cfg.topic, data)
			case <-exitchan:
				ticker.Stop()
				c.Disconnect(300)
				or := c.OptionsReader()
				sugar.Info("============================")
				sugar.Infow("disconnect!", "clientId", or.ClientID())
				return
			}
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	sugar.Infow("get os signal", "signal", <-sigChan)
	close(exitchan)
	time.Sleep(time.Second)
}
