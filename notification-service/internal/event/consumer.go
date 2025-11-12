package event

import "github.com/jackc/pgx/v5/pgxpool"

type GenericEvent struct {
	EventType   string `json:"eventType"` // e.g., order.created, payment.success
	UserID      string `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	OrderID     string `json:"orderId,omitempty"`
	Message     string `json:"message,omitempty"`
	Reservation string `json:"reservationId,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

type Consumer struct {
	Kafka       *KafkaConsumer
	Collection  *pgxpool.Pool
	ServiceName string
}

func NewConsumer(brokerURL, topic, groupID, serviceName string, collection *pgxpool.Pool) *Consumer {
	return &Consumer{
		Kafka:       NewKafkaConsumer(brokerURL, topic, groupID),
		Collection:  collection,
		ServiceName: serviceName,
	}
}
