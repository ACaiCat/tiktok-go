package ai

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var LocalToolRegistry toolRegistry

func InitLocalToolRegistry() {
	LocalToolRegistry = toolRegistry{
		tools: make(map[string]Tool),
	}

	LocalToolRegistry.RegisterTool(jwchLoginTool())
}

type toolRegistry struct {
	tools map[string]Tool
}

func (r *toolRegistry) RegisterTool(tool Tool) {
	r.tools[tool.GetName()] = tool
}

func (r *toolRegistry) ExistTool(tc openai.ToolCall) bool {
	_, ok := r.tools[tc.Function.Name]
	return ok
}

func (r *toolRegistry) CallTool(tc openai.ToolCall, callCtx ToolCallContext) (*openai.ChatCompletionMessage, error) {
	tool, ok := r.tools[tc.Function.Name]
	if !ok {
		return nil, fmt.Errorf("tool %s not found", tc.Function.Name)
	}
	callResult, err := tool.CallTool(tc, callCtx)
	if err != nil {
		return nil, err
	}
	return callResult, nil
}

func (r *toolRegistry) ListTools() []openai.Tool {
	var tools []openai.Tool

	for _, tool := range r.tools {
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.GetName(),
				Description: tool.GetDescription(),
				Parameters:  tool.GetParametersDef(),
			},
		})
	}

	return tools
}
