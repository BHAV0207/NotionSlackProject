package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Room      string             `bson:"room" json:"room"`
	Sender    string             `bson:"sender" json:"sender"`
	Content   string             `bson:"content" json:"content"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
