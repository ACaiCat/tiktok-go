package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	service "github.com/ACaiCat/tiktok-go/biz/service/user"
	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	mcp "github.com/mark3labs/mcp-go/client"
	mcptransport "github.com/mark3labs/mcp-go/client/transport"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"
)

const ChatPrompt = `
你是一个私聊场景中的AI聊天助手，只在两个用户的对话中参与。每条历史消息格式为：时间 - {{@用户ID}} 消息内容，其中 {{@AI}} 是你自己。

判断是否回复，按以下优先级执行：

1. 如果用户明确让你停止发言（如：别说话、闭嘴、noreply等类似表达），立即停止，输出 noreply
2. 如果最近消息中有人明确 @AI、向你提问、或明显是在继续与你对话，输出你的回复
3. 否则，视为用户之间的对话，输出 noreply

教务处工具：

1. 如果你的上下文没有jwch_cookie和jwch_id的信息，那么你就需要调用教务处登录
2. 你只能使用用户ID登录教务处
3. 每次会话不保存jwch_cookie和jwch_id的信息，你需要重新登录

回复要求：
- 只输出聊天内容本身，纯文本，不使用任何 Markdown 格式
- 自然简洁，符合私聊语气
- 可以用 {{@用户ID}} 称呼用户
- 不要解释判断过程，不要重复历史消息
`

type CallJwchLogin struct {
	UserID int64 `json:"user_id"`
}

func (s *ChatService) ChatAI(ctx context.Context, history string) (bool, string, error) {
	cfg := openai.DefaultConfig(config.AppConfig.AI.Key)
	cfg.BaseURL = config.AppConfig.AI.BaseURL
	client := openai.NewClientWithConfig(cfg)

	transport, err := mcptransport.NewStreamableHTTP("https://fzuhelper.west2.online/mcp")
	if err != nil {
		return false, "", err
	}

	defer transport.Close()

	mcpClient := mcp.NewClient(transport)

	defer mcpClient.Close()

	_, err = mcpClient.Initialize(ctx, mcpproto.InitializeRequest{})
	if err != nil {
		return false, "", err
	}

	listResult, err := mcpClient.ListTools(ctx, mcpproto.ListToolsRequest{})

	if err != nil {
		return false, "", err
	}

	var tools []openai.Tool

	for _, tool := range listResult.Tools {
		tools = append(tools, openai.Tool{
			Type: "",
			Function: &openai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	tools = append(tools, openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "教务处登录",
			Description: "只能使用这个工具登录，使用用户ID登录教务处，在上下文没有Cookie或账号密码时优先调用",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"user_id": map[string]interface{}{
						"type":        "integer",
						"description": "用户ID",
					},
				},
				"required": []string{"user_id"},
			},
		},
	})

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: ChatPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: history,
		},
	}

	for {
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    config.AppConfig.AI.Model,
				Messages: messages,
				Tools:    tools,
			},
		)

		log.Println(messages)

		if err != nil {
			return false, "", err
		}

		msg := resp.Choices[0].Message

		if len(msg.ToolCalls) == 0 {
			reply := msg.Content

			if reply == "noreply" {
				return false, "", nil
			}
			return true, reply, nil
		}

		for _, tc := range msg.ToolCalls {
			if tc.Function.Name == "教务处登录" {
				// 反序列化
				var param CallJwchLogin

				err := json.Unmarshal([]byte(tc.Function.Arguments), &param)
				if err != nil {
					result := &mcpproto.CallToolResult{
						Content: []mcpproto.Content{
							mcpproto.TextContent{
								Text: fmt.Sprintf("call tool error: %v", err),
							},
						},
						IsError: true,
					}
					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    extractMCPText(result),
						ToolCallID: tc.ID,
					})

					continue
				}

				id, cookie, err := service.NewUserService(ctx).GetJwchIdentifierAndCookies(param.UserID)
				if err != nil {
					result := &mcpproto.CallToolResult{
						Content: []mcpproto.Content{
							mcpproto.TextContent{
								Text: fmt.Sprintf("call tool error: %v", err),
							},
						},
						IsError: true,
					}
					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    extractMCPText(result),
						ToolCallID: tc.ID,
					})

					continue
				}

				jsonResult, _ := json.Marshal(map[string]interface{}{
					"jwch_id":     id,
					"jwch_cookie": cookie,
				})

				result := &mcpproto.CallToolResult{
					Content: []mcpproto.Content{
						mcpproto.TextContent{
							Text: string(jsonResult),
						},
					},
				}

				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    extractMCPText(result),
					ToolCallID: tc.ID,
				})

				continue
			}

			result, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
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

			messages = append(messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    extractMCPText(result),
				ToolCallID: tc.ID,
			})
		}
	}
}

func extractMCPText(result *mcpproto.CallToolResult) string {
	if result == nil {
		return ""
	}
	var texts []string
	for _, c := range result.Content {
		if tc, ok := c.(mcpproto.TextContent); ok {
			texts = append(texts, tc.Text)
		}
	}
	b, _ := json.Marshal(texts)
	return string(b)
}

func BuildAIHistory(messages []*model.ChatMessage, userAID int64, userBID int64) string {
	var history strings.Builder

	history.WriteString(fmt.Sprintf("用户A的ID: %d, 用户B的ID: %d\n", userAID, userBID))

	for _, message := range messages {
		identity := "{{@AI}}"
		if !message.IsAi {
			identity = fmt.Sprintf("{{@%d}}", message.SenderID)
		}

		history.WriteString(message.CreatedAt.Format("2006-01-02 15:04:05"))
		history.WriteString(" - ")
		history.WriteString(identity)
		history.WriteString(": ")
		history.WriteString(message.Content)
		history.WriteByte('\n')
	}

	return history.String()
}
