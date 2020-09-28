package kafkaesque

import (
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Consumer holds the observables and channels in one struct.
type Consumer struct {
	Config       *Config
	Topics       []string
	eventChannel chan rxgo.Item
	infoChannel  chan rxgo.Item
	Infos        rxgo.Observable
	Events       rxgo.Observable
}

func (c *Consumer) consumeToChannel(consumer *kafka.Consumer) {
	defer consumer.Close()
	for {
		//TODO: Add break point
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			c.eventChannel <- rxgo.Of(e)
		case kafka.PartitionEOF:
			c.infoChannel <- rxgo.Of(e)
		case kafka.Error:
			c.infoChannel <- rxgo.Of(e)
			break
		default:
			c.infoChannel <- rxgo.Of(e)
		}
	}
}

// Open start the consumer.
func (c *Consumer) Open() error {
	if len(c.Topics) == 0 {
		return fmt.Errorf("Creating a consumer without topics is not permitted")
	}
	c.eventChannel = make(chan rxgo.Item)
	c.infoChannel = make(chan rxgo.Item)
	consumer, err := kafka.NewConsumer(c.Config.Map())
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

// Close stop the consumer.
func (c *Consumer) Close() {
	//TODO: implement close channel
	return
}

// NewConsumer creates a Consumer via config and topics.
func NewConsumer(config *Config, topics []string) (*Consumer, error) {
	consumer := &Consumer{
		Config: config,
		Topics: topics,
	}
	return consumer, consumer.Open()
}

// Event mapping of confluent Message
type Event kafka.Message