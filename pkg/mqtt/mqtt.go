package mqtt

import (
	"fmt"
	"sync"
	"time"

	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt struct {
	Addr     string
	UserName string
	PassWord string
	Clientid string
}

var (
	c MQTT.Client
)

var onLost MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	or := client.OptionsReader()
	logger.Logger().Infow("mqtt connection lost", "clientId", or.ClientID(), "error", err)
}

var df = func(client MQTT.Client, msg MQTT.Message) {
	logger.Logger().Info(msg.Topic())
	logger.Logger().Info(string(msg.Payload()))
}

func Init(opt *option.Options) error {
	opts := MQTT.NewClientOptions().AddBroker(opt.Mqtt.Addr)
	if opt.Mqtt.Clientid != "" {
		opts.SetClientID(opt.Mqtt.Clientid)
	}
	//opts.SetDefaultPublishHandler(f)
	opts.SetUsername(opt.Mqtt.Username)
	opts.SetPassword(opt.Mqtt.Password)
	opts.SetCleanSession(opt.Mqtt.CleanSession)
	opts.SetConnectionLostHandler(onLost)

	//create and start a client using the above ClientOptions
	c = MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func Stop(time uint) {
	c.Disconnect(time)
}

func Sub(topic string, qos byte) (chan []byte, error) {
	mc := make(chan []byte, 1)
	f := func(client MQTT.Client, msg MQTT.Message) {
		mc <- msg.Payload()
	}
	token := c.Subscribe(topic, qos, f)
	if token.WaitTimeout(3*time.Second) == false {
		return nil, fmt.Errorf("mqtt subscribe timeout")
	}
	return mc, token.Error()
}

func Pub(wg *sync.WaitGroup, topic string, qos byte, mc chan []byte) {
	go func() {
		wg.Add(1)
		for v := range mc {
			token := c.Publish(topic, qos, false, v)
			token.Wait()
		}
		wg.Done()
	}()
}

func PubOne(topic string, qos byte, msg []byte) error {
	token := c.Publish(topic, qos, false, msg)
	token.Wait()
	return token.Error()
}
