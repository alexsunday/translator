package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConf(t *testing.T) {
	text := `
LLM openai

BASE https://api.deepseek.com
KEY sk-1234567890abcdef1234567890abcdef
MODEL deepseek-chat-1.0	`
	conf, err := ParseFile(strings.NewReader(text))
	require.Nil(t, err)
	print(conf.String())
}
