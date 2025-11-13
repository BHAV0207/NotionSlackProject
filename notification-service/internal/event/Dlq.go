package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func SendToDlq(brokerURL, dlqTopic string, event any) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    dlqTopic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})

	if err != nil {
		return fmt.Errorf("failed to write to DLQ: %v", err)
	}
	fmt.Printf("💀 Sent to DLQ [%s]: %+v\n", dlqTopic, event)
	return nil
}
