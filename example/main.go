package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/reactivex/rxgo/v2"

	"github.com/PlantUp/kafkaesque"
)

const (
	// ConsumerMode consume a topic
	ConsumerMode = "consumer"
	// ProducerMode produce from Observable
	ProducerMode = "producer"
)

// BookRequest an example type
type BookRequest struct {
	Book   string `json:"book"`
	Person string `json:"person"`
}

func main() {
	mode := flag.String("mode", "consumer", "either in producer or consumer mode")
	topics := flag.String("topics", "test", "Comma-separeted list of topics")
	broker := flag.String("broker", "localhost:9092", "Comma-seperated list of brokers")
	name := flag.String("name", "test", "ClientID")
	flag.Parse()
	brokerList := strings.Split(*broker, ",")
	topicList := strings.Split(*topics, ",")
	switch *mode {
	case ConsumerMode:
		consumeTopic(
			brokerList,
			topicList,
			*name,
		)
	case ProducerMode:
		produceTopic(
			brokerList,
			topicList,
			*name,
		)
	default:
		fmt.Println("Error no such mode")
		os.Exit(2)
	}
}

func consumeTopic(brokers []string, topics []string, clientID string) {
	testConsumer, err := kafkaesque.NewConsumer(&kafkaesque.Config{
		AutoOffsetReset:  kafkaesque.AutoOffsetReset.Latest,
		BootstrapServers: brokers,
		IsolationLevel:   kafkaesque.IsolationLevel.ReadCommitted,
		MessageMaxBytes:  1048588,
		ClientID:         clientID,
		GroupID:          clientID + "_group",
		SecurityProtocol: kafkaesque.SecurityProtocol.Plaintext,
	},
		topics...,
	)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Showing events for \u001b[34m%s\u001b[0m on \u001b[31m%s\u001b[0m:", topics, brokers))
	bookRequests := testConsumer.Events.
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			var request BookRequest
			err := json.Unmarshal(kafkaesque.ValueToBytes(item), &request)
			return request, err
		})
	defer testConsumer.Close()
	for event := range bookRequests.Observe() {
		bookRequest := event.V.(BookRequest)
		log.Println(fmt.Sprintf(
			"Book:\t%s\tPerson:\t%s",
			bookRequest.Book,
			bookRequest.Person,
		))
	}
}

func produceTopic(brokers []string, topics []string, clientID string) {
	events := rxgo.Defer(
		[]rxgo.Producer{
			func(_ context.Context, ch chan<- rxgo.Item) {
				for name, title := range map[string]string{
					"Pamela Osborne":  "10 Awesome Ways to Photograph Blue Bottles",
					"Inayah Brown":    "21 Myths About Blue bottles Debunked",
					"Ikrah Blair":     "How to Make Your Own Vast Hat for less than £5",
					"Aurora Stafford": "An analysis of handsome jugs",
					"Hafsa Scott":     "Why Are My Elbows Growing?",
				} {
					ch <- rxgo.Of(BookRequest{Person: name, Book: title})
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				}
			},
		}).Map(
		func(_ context.Context, item interface{}) (interface{}, error) {
			return func(br BookRequest) ([]byte, error) {
				return json.Marshal(&br)
			}(item.(BookRequest))
		},
	)

	testProducer, err := kafkaesque.NewProducer(
		&kafkaesque.Config{
			AutoOffsetReset:  kafkaesque.AutoOffsetReset.Latest,
			BootstrapServers: brokers,
			IsolationLevel:   kafkaesque.IsolationLevel.ReadCommitted,
			MessageMaxBytes:  1048588,
			ClientID:         clientID,
			GroupID:          clientID + "_group",
			SecurityProtocol: kafkaesque.SecurityProtocol.Plaintext,
		},
		topics[0],
		[]byte("BookRequests"),
		events,
	)

	defer testProducer.Close()
	if err != nil {
		panic(err)
	}
	for info := range testProducer.GetInfos().Observe() {
		fmt.Println(info)
	}
}
