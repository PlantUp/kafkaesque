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
	closeChannel chan bool
	Infos        rxgo.Observable
	Events       rxgo.Observable
}

// InfoEvent hold error and message
type InfoEvent struct {
	Error   error
	Message Event
}

func (c *Consumer) consumeToChannel(consumer *kafka.Consumer) {
	defer consumer.Close()
	for {
		select {
		case <-c.closeChannel:
			return
		default:
			ev, err := consumer.ReadMessage(-1)
			if err != nil {
				c.infoChannel <- rxgo.Of(
					&InfoEvent{
						Error:   err,
						Message: ev,
					},
				)
				continue
			}
			c.eventChannel <- rxgo.Of(Event(ev))
		}
	}
}

// Open start the consumer.
func (c *Consumer) Open() (*Consumer, error) {
	if len(c.Topics) == 0 {
		return nil, fmt.Errorf("Creating a consumer without topics is not permitted")
	}
	consumer, err := kafka.NewConsumer(c.Config.Map())
	if err != nil {
		return nil, err
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
	return c, nil
}

// Close stop the consumer.
func (c *Consumer) Close() {
	c.closeChannel <- true
}

// NewConsumer creates a Consumer via config and topics.
func NewConsumer(config *Config, topics ...string) (*Consumer, error) {
	consumer := &Consumer{
		Config: config,
		Topics: topics,
		eventChannel: make(chan rxgo.Item),
		infoChannel: make(chan rxgo.Item),
	}
	return consumer.Open()
}
