package kafkaesque

import (
	"fmt"
)

func TestConsume() {
	config := &Config{
		AutoOffsetReset:  AutoOffsetReset.Latest,
		BootstrapServers: []string{"localhost:9092"},
		ClientId:         "testConsumer",
		IsolationLevel:   IsolationLevel.ReadCommited,
		SecurityProtocol: SecurityProtocol.Plaintext,
		GroupId:          "1",
	}
	consumer, err := NewConsumer(config, []string{"test"})
	if err != nil {
		panic(err)
	}
	for event := range consumer.Events.Observe() {
		fmt.Println("%v", event)
	}
}
