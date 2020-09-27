package kafkaesque

import(
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/reactivex/rxgo" //"github.com/reactivex/rxgo/v2"
	"fmt"
	"log"
)

type Consumer struct {
	Config			*Config,
	Topics			[]string,
	eventChannel 	chan interface{},
	infoChannel 	chan interface{},
	Infos			*rxgo.Observable,
	Events			*rxgo.Observable,
}

func (c *Consumer) consumeToChannel(consumer *kafka.Consumer) {
	defer consumer.Close()
	for {
		//TODO: Add break point
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			c.eventChannel <- e
		case kafka.PartitionEOF:
			c.infoChannel  <- e
		case kafka.Error:
			c.infoChannel  <- e
			break
		default:
			c.infoChannel  <- e
		}
	}
}

func (c *Consumer) Open() err {
	if len(topics) == 0 {
		return fmt.Errorf("Creating a consumer without topics is not permitted.")
	}
	c.eventChannel = make(chan interface{})
	c.infoChannel  = make(chan interface{})
	consumer, err := kafka.NewConsumer(c.Config)
	if err != nil {
		return err
	}
	consumer.SubscribeTopics(c.Topics, nil)
	go c.consumeToChannel(consumer)
	c.Events = rxgo.FromEventSource(
		c.eventChannel,
		rxgo.WithBackPressureStrategy(rxgo.Drop),
	)
	c.Infos = rxgo.FromEventSource(
		c.infoChannel,
		rxgo.WithBackPressureStrategy(rxgo.Drop),
	)
	return nil
}

func (c *Consumer) Close() {
	//TODO: implement close channel
	return
}

func NewConsumer(config *Config, Topics []string) (*Consumer, err) {
	consumer := &Consumer{
		Config: config,
		Topics: topics,
	}
	return consumer, consumer.Open()
}