package ai

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"

	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

type LocalTool[I, O any] struct {
	Name          string
	Description   string
	ParametersDef any
	Authorize     func(callCtx ToolCallContext, params I) error
	Func          func(params I) (O, error)
}

func (l LocalTool[I, O]) GetName() string {
	return l.Name
}

func (l LocalTool[I, O]) GetDescription() string {
	return l.Description
}

func (l LocalTool[I, O]) GetParametersDef() map[string]any {
	r := jsonschema.Reflector{
		DoNotReference: true,
	}

	var input I
	schema := r.Reflect(input)

	data, _ := json.Marshal(schema)

	var result map[string]any
	_ = json.Unmarshal(data, &result)

	return result
}

func (l LocalTool[I, O]) CallTool(tc openai.ToolCall, callCtx ToolCallContext) (*openai.ChatCompletionMessage, error) {
	var input I

	err := json.Unmarshal([]byte(tc.Function.Arguments), &input)
	if err != nil {
		return nil, err
	}

	if l.Authorize != nil {
		err = l.Authorize(callCtx, input)
		if err != nil {
			if _, ok := any(err).(errno.ErrNo); !ok {
				err = errno.AuthErr.WithError(err)
			}

			return nil, err
		}
	}

	var callToolResult mcpproto.CallToolResult

	result, err := l.Func(input)
	if err != nil {
		callToolResult = mcpproto.CallToolResult{
			Content: []mcpproto.Content{
				mcpproto.TextContent{
					Text: fmt.Sprintf("call tool error: %v", err),
				},
			},
			IsError: true,
		}
		log.Printf("CallTool failed, tool=%s, err=%v\n", l.GetName(), err)
	} else {
		jsonResult, err := json.Marshal(result)

		if err != nil {
			return nil, err
		}

		callToolResult = mcpproto.CallToolResult{
			Content: []mcpproto.Content{
				mcpproto.TextContent{
					Text: string(jsonResult),
				},
			},
			IsError: false,
		}
	}

	callToolResultJSON, err := json.Marshal(callToolResult)

	if err != nil {
		return nil, err
	}

	msg := openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    string(callToolResultJSON),
		ToolCallID: tc.ID,
	}

	return &msg, nil
}
