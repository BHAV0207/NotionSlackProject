package kafka

import (
	"context"
	"log"

	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Kafka *KafkaProducer
}

func NewProducer(brokerURL, topic string) *Producer {
	return &Producer{
		Kafka: KafkaWriter(brokerURL, topic),
	}
}

func (p *Producer) Publish(event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: data,
	}

	err = p.Kafka.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("❌ Kafka publish failed: %v", err)
		return err
	}

	log.Printf("✅ Published event to Kafka: %s", string(data))
	return nil
}


