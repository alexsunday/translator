package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LangChainModel struct {
	llm llms.Model
}

func NewLangChainLLM(cfg *Config) (LLModel, error) {
	openaiLlm, err := openai.New(
		openai.WithBaseURL(cfg.Base),
		openai.WithModel(cfg.Model),
		openai.WithToken(cfg.Key),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI LLM: %w", err)
	}
	return &LangChainModel{
		llm: openaiLlm,
	}, nil
}

func (m *LangChainModel) GenerateContent(ctx context.Context, system string, text string, ch chan<- *translateMessage) error {
	var err error

	_, err = m.llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, system),
			llms.TextParts(llms.ChatMessageTypeHuman, text),
		},
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			var frame = &translateMessage{
				text: string(chunk),
				cmd:  translateCmd,
			}
			select {
			case <-ctx.Done():
				return nil
			case ch <- frame:
				return nil
			}
		}),
	)

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		if errors.Is(err, context.Canceled) {
			return nil
		}
		ch <- &translateMessage{
			text: err.Error(),
			cmd:  errorCmd,
		}

		return fmt.Errorf("failed to generate content: %w", err)
	}

	println(finishedCmd)
	ch <- &translateMessage{
		text: "Translation completed",
		cmd:  finishedCmd,
	}
	return nil
}
