package option

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Logging struct {
	Level string
	Path  string
}

type DataSource struct {
	Driver   string
	Host     string
	Port     string
	UserName string
	PassWord string
	DbName   string
	Sslmode  string
	LogLevel string `yaml:"log-level"`
}

type Kafka struct {
	Addrs    []string
	Clientid string
	Groupid  string
	Device   string
	Topic    struct {
		Device string
	}
}

type Mqtt struct {
	Addr         string
	Clientid     string
	Username     string
	Password     string
	CleanSession bool `yaml:"clean-session"`
}

type Minio struct {
	Endpoint        string
	Accesskeyid     string
	Secretaccesskey string
	Usessl          bool
	Bucketname      string
	Location        string
	Prefix          string
}

type MqttTopic struct {
	Topic    string
	Qos      byte
	Fmtstr   string
	Interval int
}

type Camera struct {
	Ip       string
	Username string
	Password string
}

type Options struct {
	viper *viper.Viper

	Logging
	DataSource
	Kafka
	Mqtt
	Minio
	MqttTopics []MqttTopic `yaml:"mqtt-topics"`
	Camera
}

func New() *Options {
	opt := &Options{
		viper: viper.New(),
	}
	return opt
}

func (opt *Options) Parse() error {
	opt.viper.SetConfigFile("app.yml")
	opt.viper.SetConfigType("yaml")
	err := opt.viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file app.yml failed: %v", err)
	}

	err = opt.viper.Unmarshal(opt, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	})
	if err != nil {
		return fmt.Errorf("yaml file unmarshal failed, please make sure you provide valid yaml file, %v", err)
	}

	return nil
}
