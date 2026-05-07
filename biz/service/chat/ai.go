package service

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/ai"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/sashabaranov/go-openai"
)

func (s *ChatService) replyWithAI(userID int64, receiverID int64) {
	messages, err := s.getChatHistory(userID, receiverID, 10, 0)
	if err != nil {
		log.Println("failed to get message history:", err)
		return
	}

	// 反转使最新的消息在上下文底部
	slices.Reverse(messages)
	history := buildAIHistory(messages, userID, receiverID)

	reply, content, err := chatWithAI(s.ctx, history, userID, receiverID)
	if err != nil {
		log.Println("failed to AI message:", err)
		return
	}
	if !reply {
		return
	}

	now := time.Now().UnixMilli()
	receiverMessage := &ws.ChatMessage{
		SenderID:   userID,
		ReceiverID: receiverID,
		IsAI:       true,
		Content:    content,
		Timestamp:  now,
	}
	receiverOnline := s.sendMessageToUser(receiverID, ws.MessageTypeChat, receiverMessage, "failed to forward message to receiver")
	if !s.saveChatMessage(userID, receiverID, content, receiverOnline, true) {
		return
	}

	senderMessage := &ws.ChatMessage{
		SenderID:   receiverID,
		ReceiverID: userID,
		IsAI:       true,
		Content:    content,
		Timestamp:  now,
	}
	senderOnline := s.sendMessageToUser(userID, ws.MessageTypeChat, senderMessage, "failed to forward message to sender")
	_ = s.saveChatMessage(receiverID, userID, content, senderOnline, true)
}

type CallJwchLogin struct {
	UserID int64 `json:"user_id"`
}

func chatWithAI(ctx context.Context, history string, userAID int64, userBID int64) (bool, string, error) {
	cfg := openai.DefaultConfig(config.AppConfig.AI.Key)
	cfg.BaseURL = config.AppConfig.AI.BaseURL
	client := openai.NewClientWithConfig(cfg)

	fuuMcp, err := ai.NewFuuMCPClient(ctx)
	if err != nil {
		return false, "", err
	}

	tools, err := fuuMcp.ListTools()

	if err != nil {
		return false, "", err
	}

	tools = append(tools, ai.LocalToolRegistry.ListTools()...)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constants.Prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: history,
		},
	}

	reply, err := agentLoop(ctx, client, fuuMcp, messages, tools, ai.NewToolCallContext(userAID, userBID))

	if err != nil {
		return false, "", err
	}

	if reply == constants.NoReplySignal {
		return false, "", nil
	}

	return true, reply, nil

}

func agentLoop(ctx context.Context, client *openai.Client, fuuMCP *ai.FuuMCP,
	messages []openai.ChatCompletionMessage, tools []openai.Tool, localToolCallCtx ai.ToolCallContext) (string, error) {
	for {
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    config.AppConfig.AI.Model,
				Messages: messages,
				Tools:    tools,
			},
		)

		if err != nil {
			return "", err
		}

		msg := resp.Choices[0].Message
		messages = append(messages, msg)

		if len(msg.ToolCalls) == 0 {
			return msg.Content, nil
		}

		for _, tc := range msg.ToolCalls {
			var callMsg *openai.ChatCompletionMessage
			if ai.LocalToolRegistry.ExistTool(tc) {
				callMsg, err = ai.LocalToolRegistry.CallTool(tc, localToolCallCtx)
				if err != nil {
					return "", err
				}
			} else {
				callMsg, err = fuuMCP.CallTool(tc)
				if err != nil {
					return "", err
				}
			}

			messages = append(messages, *callMsg)
		}
	}
}

func buildAIHistory(messages []*model.ChatMessage, userAID int64, userBID int64) string {
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
