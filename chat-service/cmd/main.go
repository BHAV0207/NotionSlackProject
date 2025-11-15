package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/chat-service/internal/repository"
	"github.com/BHAV0207/chat-service/internal/websockets"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	hub    = websockets.NewHub()
	client *mongo.Client
)



func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI not set")
	}

	client = db.ConnectMongo(uri)
	go hub.Run()


	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWS(hub, client, w, r)
	})

	fmt.Println("🚀 Chat server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
