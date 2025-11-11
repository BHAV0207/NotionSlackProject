package kafka

import (
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func KafkaWriter(brokerURL, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokerURL), // Connect to the Kafka broker
			Topic:        topic,                // The topic to publish to
			Balancer:     &kafka.LeastBytes{},  // Choose the partition with least data
			RequiredAcks: kafka.RequireAll,     // Wait for all replicas to confirm
		},
	}
}
