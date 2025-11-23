package services

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMService struct {
	llm       llms.LLM
	provider  string // open ai or any other
	modelName string
}

func NewLLMService() (*LLMService, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")

	// Try OpenAI first if key is present
	if openAIKey != "" {
		modelName := os.Getenv("OPENAI_MODEL")
		if modelName == "" {
			modelName = "gpt-4o-mini"
		}

		llm, err := openai.New(
			openai.WithToken(openAIKey),
			openai.WithModel(modelName),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init OpenAI LLM: %w", err)
		}

		return &LLMService{
			llm:       llm,
			provider:  "openai",
			modelName: modelName,
		}, nil
	}

	// Fallback to Ollama
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "llama3"
	}

	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}

	llm, err := ollama.New(
		ollama.WithModel(ollamaModel),
		ollama.WithServerURL(ollamaHost),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to init Ollama LLM: %w", err)
	}

	return &LLMService{
		llm:       llm,
		provider:  "ollama",
		modelName: ollamaModel,
	}, nil
}

func (s *LLMService) Provider() string {
	return s.provider
}

func (s *LLMService) ModelName() string {
	return s.modelName
}

func (s *LLMService) GenerateText(ctx context.Context, prompt string) (string, error) {
	if s.llm == nil {
		return "", fmt.Errorf("LLM not initialized")
	}
	

	resp, err := s.llm.Call(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp, nil
}

func (s *LLMService) Summarize(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf(`
	You are a helpful assistant that summarizes text for users of a real-time collaboration platform.
	
	Summarize the following content in a clear, concise way.
	- Keep it under 5 bullet points.
	- Preserve all key decisions and action items.
	
	Content:
	%s
	
	Summary:
	`, text)

	return s.GenerateText(ctx, prompt)
}

// Chat runs a chat-style interaction using system + user messages.
func (s *LLMService) Chat(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	if s.llm == nil {
		return "", fmt.Errorf("LLM not initialized")
	}

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userMessage),
	}

	// GenerateContent returns *llms.ContentResponse
	resp, err := s.llm.GenerateContent(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM chat failed: %w", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("LLM returned no choices")
	}

	choice := resp.Choices[0]

	if choice.FuncCall != nil {
		return "", fmt.Errorf("model requested a function call, not plain text")
	}

	// The actual text content
	return choice.Content, nil
}

func (s *LLMService) Suggest(ctx context.Context, contextText string) (string, error) {
	systemPrompt := `
You are an AI assistant inside a real-time collaboration & chat platform.
Given the current conversation or document content, suggest helpful replies,
follow-up questions, or ideas. Be specific and actionable.
Return the output as a small numbered list.
`
	userMessage := fmt.Sprintf(`
Here is the current message or content context:

%s

Generate suggested replies or follow-up actions.
`, contextText)

	return s.Chat(ctx, systemPrompt, userMessage)
}
