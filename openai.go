package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAILLm struct {
	Model  string
	client *openai.Client
}

func NewOpenAiLLM(baseURL, model, token string) *OpenAILLm {
	cfg := openai.DefaultConfig(token)
	cfg.BaseURL = baseURL

	return &OpenAILLm{
		Model:  model,
		client: openai.NewClientWithConfig(cfg),
	}
}

func (llm *OpenAILLm) GenerateContent(ctx context.Context, system string, text string, ch chan<- *translateMessage) error {
	req := openai.ChatCompletionRequest{
		Model: llm.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: system,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		Stream: true,
	}
	stream, err := llm.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create chat completion stream: %w", err)
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			ch <- &translateMessage{
				cmd: finishedCmd,
			}

			return nil
		}

		if err != nil {
			ch <- &translateMessage{
				text: err.Error(),
				cmd:  errorCmd,
			}
			return nil
		}

		content := response.Choices[0].Delta.Content
		if content != "" {
			ch <- &translateMessage{
				text: content,
				cmd:  translateCmd,
			}
		}
	}
}
