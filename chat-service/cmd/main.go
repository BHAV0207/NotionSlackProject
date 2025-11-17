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
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ .env not found — using environment variables")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		fmt.Println("⚠️ MONGO_URI not set (skipping Mongo connection)")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client = db.ConnectMongo(uri)

	// init history route so /history is registered
	initHistoryRoutes(client)

	// create redis broker (use env or defaults)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redisPass := os.Getenv("REDIS_PASS") // optional
	broker := websockets.NewRedisBroker(redisAddr, redisPass)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWS(hub, broker, client, w, r) // pass both broker and db client
	})

	fmt.Printf("🚀 Chat server running on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
