package event

import "github.com/segmentio/kafka-go"

type KafkaConsumer struct {
	Reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer instance
func NewKafkaConsumer(brokerURL, topic, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerURL},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}
