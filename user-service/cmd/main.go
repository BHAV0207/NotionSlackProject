package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BHAV0207/user-service/internal/handler"
	"github.com/BHAV0207/user-service/internal/kafka"
	"github.com/BHAV0207/user-service/internal/middleware"
	"github.com/BHAV0207/user-service/internal/repository"
	"github.com/BHAV0207/user-service/pkg/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ No .env file found, using system environment variables")
	}

	URI := os.Getenv("DB_URL")
	if URI == "" {
		panic("DB_URL not set in environment")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	client := repository.ConnectDb(URI)
	defer client.Close()

	redis.InitRedis()
	// start HTTP server or consumer

	userCreatedEvent := kafka.NewProducer(broker, "user-created")
	h := &handler.UserHandler{
		DB:       client,
		Producer: userCreatedEvent,
	}

	router := mux.NewRouter()
	router.HandleFunc("/register", h.RegisterUser).Methods("POST")
	router.HandleFunc("/login", h.LoginUser).Methods("POST")
	router.HandleFunc("/{id}", h.GetUserByID).Methods("GET")
	router.Handle("/me", middleware.JWTAuth(http.HandlerFunc(h.UserInfo))).Methods("GET")

	if err := http.ListenAndServe(":"+PORT, router); err != nil {
		panic(err)
	}

}
