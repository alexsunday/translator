package main

import (
	"fmt"
	"os"
)

const (
	defaultSystemPrompt = `你是一个专业的翻译员, 精通中文与英文的翻译工作, 熟谙各种翻译技巧和方法, 能够准确地将中文翻译成英文, 并且能够理解和传达原文的意思和情感。
	你主要工作在 IT/计算机领域, 需要翻译各种技术文档、代码注释、用户手册等内容。请确保翻译的内容准确、流畅，并且符合目标语言的语法和用词习惯。
	你要自动判断用户的输入语言, 如果是中文, 则将其翻译成英文。如果是英文, 则将其翻译成中文。请注意, 你只需要翻译内容, 不需要解释或评论原文。
	不需要输出思考过程，不需要输出翻译说明。只需要输出翻译后的内容。`
	llmOllama = "ollama"
	llmOpenAI = "openai"
)

type Config struct {
	LLM    string
	Key    string
	Model  string
	Base   string
	System string
}

func (c *Config) String() string {
	return fmt.Sprintf("LLM: %s, Key: %s, Model: %s, Base: %s, System: %s",
		c.LLM, c.Key, c.Model, c.Base, c.System)
}

func fromModalFile(m *Modelfile) (*Config, error) {
	if m == nil {
		return nil, fmt.Errorf("modelfile is nil")
	}

	// llm, base, key, model 只能有一个 如果有多个报错
	var dict = make(map[string]string)
	for _, item := range m.Commands {
		if item.Name == "llm" {
			if _, exists := dict["llm"]; exists {
				return nil, fmt.Errorf("multiple llm specified")
			}
			dict["llm"] = item.Args
		} else if item.Name == "base" {
			if _, exists := dict["base"]; exists {
				return nil, fmt.Errorf("multiple base specified")
			}
			dict["base"] = item.Args
		} else if item.Name == "key" {
			if _, exists := dict["key"]; exists {
				return nil, fmt.Errorf("multiple key specified")
			}
			dict["key"] = item.Args
		} else if item.Name == "model" {
			if _, exists := dict["model"]; exists {
				return nil, fmt.Errorf("multiple model specified")
			}
			dict["model"] = item.Args
		} else if item.Name == "system" {
			if _, exists := dict["system"]; exists {
				return nil, fmt.Errorf("multiple system specified")
			}
			dict["system"] = item.Args
		}
	}

	llm := dict["llm"]
	if llm == "" {
		llm = "openai"
	}

	base := dict["base"]
	if base == "" {
		base = "https://api.openai.com/v1"
	}

	key := dict["key"]
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}

	model := dict["model"]
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	system := dict["system"]
	if system == "" {
		system = defaultSystemPrompt
	}

	return &Config{
		LLM:    llm,
		Key:    key,
		Model:  model,
		Base:   base,
		System: system,
	}, nil
}

func loadConfig() (*Config, error) {
	reader, err := os.Open("conf")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	modelFile, err := ParseFile(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config, err := fromModalFile(modelFile)
	if err != nil {
		return nil, fmt.Errorf("failed to convert modelfile to config: %w", err)
	}

	return config, nil
}
