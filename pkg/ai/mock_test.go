package ai

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type mockTool struct {
	name        string
	description string
	params      map[string]any
	callResult  *openai.ChatCompletionMessage
	callErr     error
}

func (m *mockTool) GetName() string                  { return m.name }
func (m *mockTool) GetDescription() string           { return m.description }
func (m *mockTool) GetParametersDef() map[string]any { return m.params }
func (m *mockTool) CallTool(_ openai.ToolCall, _ ToolCallContext) (*openai.ChatCompletionMessage, error) {
	return m.callResult, m.callErr
}

func newMock(name string) *mockTool {
	return &mockTool{
		name:        name,
		description: fmt.Sprintf("description of %s", name),
		params:      map[string]any{"type": "object"},
		callResult:  &openai.ChatCompletionMessage{Content: "result of " + name},
	}
}
