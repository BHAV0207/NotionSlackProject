package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(uri string) *mongo.Client {
	clientOpts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("❌ Mongo connection error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Mongo ping failed:", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")
	return client
}
