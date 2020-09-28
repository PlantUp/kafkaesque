package kafkaesque

import (
	"github.com/reactivex/rxgo/v2"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Event mapping of confluent message
type Event *kafka.Message

// ValueToString convert event to string.
func ValueToString(e interface{}) string {
	return string(ValueToBytes(e))
}

// ValueToBytes convert event to byte array
func ValueToBytes(e interface{}) []byte {
	return []byte(e.(Event).Value)
}

// ToEvent convert Item to Event
func ToEvent(e interface{}) Event {
	return e.(rxgo.Item).V.(Event)
}