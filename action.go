package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type translateMessage struct {
	text string
	cmd  string
}

const (
	translateCmd = "translate"
	finishedCmd  = "finished"
	errorCmd     = "error"
)

func llmGenerateContent(ctx context.Context, llm llms.Model, system string, text string, ch chan<- *translateMessage) error {
	var err error

	_, err = llm.GenerateContent(
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
