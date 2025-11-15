package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var MessageCollection *mongo.Collection

func ConnectDb(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ Mongo connect failed: %v", err)
	}
	Client = client
	MessageCollection = client.Database("chat").Collection("messages")
	log.Println("✅ MongoDB connected")

}
