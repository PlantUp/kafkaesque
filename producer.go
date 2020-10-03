package kafkaesque

import (
	"github.com/reactivex/rxgo/v2"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Header mapping of kafka header
type Header kafka.Header

func mapHeader(xs []Header) []kafka.Header {
	var f []kafka.Header
	for _, x := range xs {
		f = append(f, kafka.Header(x))
	}
	return f
}

// ErrorEvent will be exported when producer crashes
type ErrorEvent struct {
	Error error
}

func (ee ErrorEvent) String() string {
	return ee.Error.Error()
}

//Producer holds the observables and channels in one struct.
type Producer struct {
	Config       *Config
	Topic        string
	Key          []byte
	closeChannel chan bool
	infoChannel  chan kafka.Event
	producer     *kafka.Producer
	Events       rxgo.Observable
	headers      []kafka.Header
}

//Open starts the Producer
func (p *Producer) open() (*Producer, error) {
	var err error
	p.producer, err = kafka.NewProducer(p.Config.Map())
	if err != nil {
		return nil, err
	}
	go func(producer *Producer) {
		defer p.producer.Close()
		for ev := range producer.Events.Observe() {
			err := producer.producer.Produce(
				&kafka.Message{
					Value:          ev.V.([]byte),
					Key:            p.Key,
					TopicPartition: kafka.TopicPartition{Topic: &producer.Topic, Partition: kafka.PartitionAny},
					Headers:        p.headers,
				},
				p.infoChannel,
			)
			if err != nil {
				p.infoChannel <- kafka.Event(ErrorEvent{Error: err})
				return
			}
		}
	}(p)
	return p, nil
}

// GetInfos return infos from producer as Observable
func (p *Producer) GetInfos() rxgo.Observable {
	items := make(chan rxgo.Item)
	go func(items chan rxgo.Item, infos chan kafka.Event) {
		for info := range infos {
			items <- rxgo.Of(info)
		}
	}(items, p.infoChannel)
	return rxgo.FromChannel(items)
}

// Close the producer
func (p *Producer) Close() {
	p.producer.Close()
}

// NewProducer creates a new producer instance from a observable
func NewProducer(config *Config, topic string, key []byte, events rxgo.Observable, header ...Header) (*Producer, error) {
	producer := &Producer{
		Config:       config,
		Topic:        topic,
		Key:          key,
		Events:       events,
		closeChannel: make(chan bool),
		infoChannel:  make(chan kafka.Event),
		headers:      mapHeader(header),
	}
	return producer.open()
}
