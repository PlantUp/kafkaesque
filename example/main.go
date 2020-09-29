package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PlantUp/kafkaesque"
)

// BookRequest an example type
type BookRequest struct {
	Book   string `json:"book"`
	Person string `json:"person"`
}

func main() {
	testConsumer, err := kafkaesque.NewConsumer(&kafkaesque.Config{
		AutoOffsetReset:  kafkaesque.AutoOffsetReset.Latest,
		BootstrapServers: []string{"localhost:9092"},
		IsolationLevel:   kafkaesque.IsolationLevel.ReadCommitted,
		MessageMaxBytes:  1048588,
		ClientID:         "test_1",
		GroupID:          "test_group",
		SecurityProtocol: kafkaesque.SecurityProtocol.Plaintext,
	},
		"test",
		"test2",
	)
	if err != nil {
		panic(err)
	}
	log.Println("Showing events:")
	bookRequests := testConsumer.Events.
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			var request BookRequest
			err := json.Unmarshal(kafkaesque.ValueToBytes(item), &request)
			return request, err
		})
	defer testConsumer.Close()
	for event := range bookRequests.Observe() {
		log.Println(fmt.Sprintf("%T", event.V))
	}
}
