package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	// addrs := []string{"192.168.152.185:9092", "192.168.152.48:9093", "192.168.152.48:9094"}
	topic := "alarm.business"
	addrs := []string{"192.168.152.185:9092"}
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	admin, err := sarama.NewClusterAdmin(addrs, config)
	if err != nil {
		fmt.Println(err)
	}
	err = admin.CreateTopic(topic, &sarama.TopicDetail{NumPartitions: 1, ReplicationFactor: 1}, false)
	if err != nil {
		fmt.Println(err)
	}

	err = admin.Close()
	if err != nil {
		fmt.Println(err)
	}

	producer, err := sarama.NewSyncProducer(addrs, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	strTemplate := `{
        "ts": %d,
        "uuid": "ce64eb2ee59638bcb07faba898d1bed8",
        "deviceName": "测试雷达(勿删)",
        "deviceType": 2,
        "longitude": 108.8162665,
        "latitude": 34.17049035,
        "assetId": "5de09533a09a36b991d9c86347a08bf4",
        "roleId": 0,
        "alarmName": "超速",
        "alarmLevel": 4,
		"eventType": 10,
        "imgPath": "",
        "status": 1
      }`

	for {
		str := fmt.Sprintf(strTemplate, time.Now().Unix())
		msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(str)}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println("failed to send message: ", err)
		} else {
			fmt.Printf("message sent to partition %d at offset %d\n", partition, offset)
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
