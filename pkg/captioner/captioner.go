package captioner

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultModel = openai.GPT5
)

// Captioner is the interface for a service that can caption an image.
type Captioner interface {
	Caption(ctx context.Context, imagePath string) (string, error)
}

// Client is an interface for the OpenAI client, for testing purposes.
type Client interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// OpenAICaptioner is an implementation of Captioner that uses the OpenAI API.
type OpenAICaptioner struct {
	client Client
	model  string
}

// New creates a new OpenAICaptioner.
func New(apiKey, baseURL, model string) *OpenAICaptioner {
	config := openai.DefaultConfig(apiKey)

	if baseURL != "" {
		config.BaseURL = baseURL
	}

	if model == "" {
		model = defaultModel
	}

	return &OpenAICaptioner{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// Caption generates a caption for the given image.
func (c *OpenAICaptioner) Caption(ctx context.Context, imagePath string) (string, error) {
	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}
	return c.caption(ctx, imgBytes)
}

func (c *OpenAICaptioner) caption(ctx context.Context, imgBytes []byte) (string, error) {
	encodedImage := base64.StdEncoding.EncodeToString(imgBytes)

	mimeType := http.DetectContentType(imgBytes)

	imageURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encodedImage)

	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an image captioning assistant, you provide clear and concise captions for images. You respond only with the caption.",
			},
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: "Provide a clear and concise caption suitable for screen readers for the following image:",
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: imageURL,
						},
					},
				},
			},
		},
		MaxCompletionTokens: 300,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return resp.Choices[0].Message.Content, nil
}
