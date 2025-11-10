package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BHAV0207/user-service/internal/handler"
	"github.com/BHAV0207/user-service/internal/repository"
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

	client := repository.ConnectDb(URI)
	defer client.Close()

	h := &handler.UserHandler{
		DB: client,
	}

	router := mux.NewRouter()
	router.HandleFunc("/register", h.RegisterUser).Methods("POST")
	router.HandleFunc("/login", h.LoginUser).Methods("POST")

	if err := http.ListenAndServe(":"+PORT, router); err != nil {
		panic(err)
	}

}
