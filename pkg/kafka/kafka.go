package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
)

var (
	c            sarama.Client
	groupid      string
	syncProducer sarama.SyncProducer
)

type Kafka struct {
	Addrs    []string
	Clientid string
	Groupid  string
}

type GpHandler struct {
	mch map[string]chan []byte
}

func (GpHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h GpHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	for _, ch := range h.mch {
		close(ch)
	}
	return nil
}

func (h GpHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.mch[msg.Topic] <- msg.Value
		logger.Logger().Info(msg.Topic)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func Init(opt *option.Options) (err error) {
	syncProducer, err = sarama.NewSyncProducer(opt.Kafka.Addrs, nil)
	if err != nil {
		return
	}
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Producer.MaxMessageBytes = 20000000
	// config.Producer.Return.Successes = true
	config.ClientID = opt.Kafka.Clientid
	config.Producer.Flush.Frequency = 50
	// config.Producer.Return.Successes = false
	//config.Producer.RequiredAcks = sarama.WaitForLocal
	//config.Producer.Return.Successes = true
	//config.Producer.Return.Errors = true
	//config.Producer.Flush.Bytes = 102400
	//config.Producer.Partitioner = sarama.NewHashPartitioner
	c, err = sarama.NewClient(opt.Kafka.Addrs, config)
	groupid = opt.Kafka.Groupid
	return err
}

func Stop() error {
	syncProducer.Close()
	return c.Close()
}

func PubOne(topic string, val []byte) {
	fmt.Println("send start")
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(val),
	}
	_, _, err := syncProducer.SendMessage(msg)
	if err != nil {
		logger.Logger().Error("send msg error: ", err.Error())
		return
	}
}

func Pub(wg *sync.WaitGroup, topic string) (chan []byte, error) {
	//func Pub(ctx context.Context, topic string) (chan []byte, error) {
	logger.Logger().Info("pub kafka topic: ", topic)
	in := make(chan []byte)
	producer, err := sarama.NewAsyncProducerFromClient(c)
	if err != nil {
		return nil, fmt.Errorf("kafka pub error: %s", err)
	}

	/*
		fmt.Println("kafka topic", topic)
		go func() {
			for {
				select {
				case v := <-in:
					msg := &sarama.ProducerMessage{
						Topic: topic,
						Key:   nil,
					}
					msg.Value = sarama.ByteEncoder(v)
					producer.Input() <- msg
					fmt.Printf("time: %s,topic: %s,msg len: %d\n", time.Now(), topic, msg.Value.Length())
				case err := <-producer.Errors():
					fmt.Println("send to kafka error", err)
				case <-ctx.Done():
					producer.AsyncClose()
					fmt.Println("producer close", topic)
					return
				}
			}
		}()
	*/

	go func() {
		wg.Add(1)
		for err := range producer.Errors() {
			logger.Logger().Error("send to kafka error: ", err.Error())
		}
		logger.Logger().Info(topic, " error close")
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		for v := range in {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   nil,
				Value: sarama.ByteEncoder(v),
			}
			producer.Input() <- msg
			logger.Logger().Infow("send msg success", "topic", topic, "msg_len", msg.Value.Length())
		}
		producer.AsyncClose()
		logger.Logger().Info(topic, " producer close")
		wg.Done()
	}()
	return in, nil
}

func Sub(ctx context.Context, mch map[string]chan []byte) error {
	group, err := sarama.NewConsumerGroupFromClient(groupid, c)
	if err != nil {
		return err
	}

	// Track errors
	go func() {
		for err := range group.Errors() {
			logger.Logger().Error("kafka consume group error: ", err.Error())
		}
		logger.Logger().Info("kafka consume group err done")
	}()

	// Iterate over consumer sessions.
	//ctx := context.Background()
	go func() {
		topics := make([]string, len(mch))
		i := 0
		for k := range mch {
			topics[i] = k
			i++
			logger.Logger().Info("sub kafka topic: ", k)
		}
		handler := GpHandler{mch: mch}

		for {
			err = group.Consume(ctx, topics, handler)
			if err != nil {
				return
			}
			logger.Logger().Info("group done")
		}
	}()

	go func() {
		<-ctx.Done()
		err := group.Close()
		if err != nil {
			logger.Logger().Error("groupconsumer close error: ", err.Error())
		}
	}()

	return nil
}
