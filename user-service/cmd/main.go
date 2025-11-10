package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	URI := os.Getenv("DB_URL")
	if URI == "" {
		panic("DB_URL not set in environment")
	}

	PORT:= os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	
}