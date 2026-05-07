package ai

import (
	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	GetName() string
	GetDescription() string
	CallTool(tc openai.ToolCall, callCtx ToolCallContext) (*openai.ChatCompletionMessage, error)
	GetParametersDef() map[string]any
}
