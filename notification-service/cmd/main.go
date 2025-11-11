package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BHAV0207/notification-service/internal/event"
	"github.com/BHAV0207/notification-service/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ No .env file found, using system environment variables")
	}

	URI := os.Getenv("DB_URI")
	if URI == "" {
		panic("DB_URI not set in environment")
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "9000"
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	client := repository.ConnectDb(URI)
	defer client.Close()

	userConsumer := event.NewConsumer(
		broker,
		"user-created",
		"notif-user-group",
		"notification-service",
		client,
	)

	// Run both consumers concurrently
	// go orderConsumer.StartConsuming()
	go userConsumer.StartConsuming()

	server := &http.Server{
		Addr:         ":" + PORT,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("🔔 Notification Service running on http://localhost:%s\n", PORT)
	log.Fatal(server.ListenAndServe())

}
