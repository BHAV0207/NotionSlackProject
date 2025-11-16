package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/BHAV0207/chat-service/internal/repository"
	"github.com/BHAV0207/chat-service/internal/websockets"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	hub    = websockets.NewHub()
	client *mongo.Client
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("fatal")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI not set")
	}

	client = db.ConnectMongo(uri) // import internal/db package

	// init history route so /history is registered
	initHistoryRoutes(client)

	// create redis broker (use env or defaults)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPass := os.Getenv("REDIS_PASS") // optional
	broker := websockets.NewRedisBroker(redisAddr, redisPass)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWS(hub, broker, client, w, r) // pass both broker and db client
	})

	fmt.Println("🚀 Chat server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
