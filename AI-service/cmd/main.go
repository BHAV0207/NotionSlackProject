package main

import (
	"log"

	"github.com/BHAV0207/AI-service/internal/handler"
	"github.com/BHAV0207/AI-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// init LLM
	llmService, err := services.NewLLMService()
	if err != nil {
		log.Fatalf("Failed to init LLM: %v", err)
	}

	log.Println("AI Provider:", llmService.Provider())
	log.Println("Model:", llmService.ModelName())

	aiHandler := handler.NewAiHandler(llmService)

	app.Post("/summarize", aiHandler.Summarize)
	app.Post("/suggest", aiHandler.Suggest)
	
	log.Println("AI Service running on :8005")
	log.Fatal(app.Listen(":8005"))
}
