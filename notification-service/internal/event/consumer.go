package event

import "github.com/jackc/pgx/v5/pgxpool"

type GenericEvent struct {
	EventType   string `json:"eventType"` // e.g., order.created, payment.success
	UserID      string `json:"userId"`
	OrderID     string `json:"orderId"`
	Message     string `json:"message"`
	Reservation string `json:"reservationId,omitempty"`
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
