package cmd

import (
	"fmt"
	"os"

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

	
}
