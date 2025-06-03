package main

import "context"

type translateMessage struct {
	text string
	cmd  string
}

const (
	translateCmd = "translate"
	finishedCmd  = "finished"
	errorCmd     = "error"
)

type LLModel interface {
	GenerateContent(ctx context.Context, system string, text string, ch chan<- *translateMessage) error
}
