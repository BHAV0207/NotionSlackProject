package handler

import (
	"context"

	"github.com/BHAV0207/AI-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	LLM *services.LLMService
}

func NewAiHandler(llm *services.LLMService) *AIHandler {
	return &AIHandler{LLM: llm}
}

//  in this code we are not using the http framework we are using the fiber framework of the go lang so this might look different that the http one 

func (h *AIHandler) Summarize(c *fiber.Ctx) error {
	var body struct {
		Text string `json:"text"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := h.LLM.Summarize(context.Background(), body.Text)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"summary": res,
	})
}

func (h *AIHandler) Suggest(c *fiber.Ctx) error {
	var body struct {
		Context string `json:"context"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := h.LLM.Suggest(context.Background(), body.Context)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"suggestions": res,
	})
}
