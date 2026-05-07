package ai

import (
	"context"
	"encoding/json"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/client"
	mcptransport "github.com/mark3labs/mcp-go/client/transport"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"
)

type FuuMCP struct {
	client *mcp.Client
	ctx    context.Context
}

func NewFuuMCPClient(ctx context.Context) (*FuuMCP, error) {
	transport, err := mcptransport.NewStreamableHTTP("https://fzuhelper.west2.online/mcp")
	if err != nil {
		return nil, err
	}
	mcpClient := mcp.NewClient(transport)
	_, err = mcpClient.Initialize(ctx, mcpproto.InitializeRequest{})
	if err != nil {
		return nil, err
	}

	return &FuuMCP{client: mcpClient, ctx: ctx}, nil
}

func (f *FuuMCP) ListTools() ([]openai.Tool, error) {
	listResult, err := f.client.ListTools(f.ctx, mcpproto.ListToolsRequest{})

	if err != nil {
		return nil, err
	}

	var tools []openai.Tool

	for _, tool := range listResult.Tools {
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	return tools, nil
}

func (f *FuuMCP) CallTool(tc openai.ToolCall) (*openai.ChatCompletionMessage, error) {
	result, err := f.client.CallTool(f.ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name:      tc.Function.Name,
			Arguments: json.RawMessage(tc.Function.Arguments),
		},
	})
	if err != nil {
		result = &mcpproto.CallToolResult{
			Content: []mcpproto.Content{
				mcpproto.TextContent{
					Text: fmt.Sprintf("call tool error: %v", err),
				},
			},
			IsError: true,
		}
	}

	jsonResult, err := result.MarshalJSON()

	if err != nil {
		return nil, err
	}

	msg := openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    string(jsonResult),
		ToolCallID: tc.ID,
	}

	return &msg, nil
}
