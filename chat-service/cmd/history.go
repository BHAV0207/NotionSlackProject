package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/BHAV0207/chat-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client // set from main

func historyHandler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	if room == "" {
		http.Error(w, "room required", http.StatusBadRequest)
		return
	}

	limitQ := r.URL.Query().Get("limit")
	limit := int64(50)
	if limitQ != "" {
		if v, err := strconv.Atoi(limitQ); err == nil && v > 0 {
			limit = int64(v)
		}
	}

	// optional: before timestamp (ISO8601) for pagination
	before := r.URL.Query().Get("before")
	filter := bson.M{"room_id": room}
	if before != "" {
		if t, err := time.Parse(time.RFC3339, before); err == nil {
			filter["timestamp"] = bson.M{"$lt": t}
		}
	}

	collection := mongoClient.Database("chatdb").Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // newest first
	findOpts.SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []models.Message
	if err := cursor.All(ctx, &results); err != nil {
		http.Error(w, "decode error", http.StatusInternalServerError)
		return
	}

	// Return in ascending order (oldest -> newest) for client-friendly display
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func initHistoryRoutes(mongoCli *mongo.Client) {
	mongoClient = mongoCli
	http.HandleFunc("/history", historyHandler)
}
