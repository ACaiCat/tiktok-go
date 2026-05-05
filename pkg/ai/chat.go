package ai

import (
	"context"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/sashabaranov/go-openai"
)

const ChatPrompt = `
你是一个只在私聊场景中发言的AI聊天助手。你需要根据历史消息判断当前是否需要你回复。

规则：

1. 历史消息中的每一行都形如“@名字: 消息内容”。其中@AI是你自己的发言标识，其他@用户名表示对应用户。
2. 用户可能有两个人，你需要通过不同的@用户名区分不同用户。
3. 你只需要识别这些称呼代表的发言身份，不要在回复中解释这些称呼。
4. 必须重点关注最后几条聊天内容，优先根据最近的上下文判断是否需要回复。
5. 必须遵循用户对你的明确命令。如果用户让你不要说话、别回、闭嘴、停下、安静，或表达类似意思，你就不要回复，直接输出：noreply
6. 如果当前消息没有提到你、没有点名@AI、没有向你提问，也没有明显是在继续与你对话，则默认是在和另一个用户说话，你不要回复，直接输出：noreply
7. 如果当前消息不需要你回复，直接输出：noreply
8. 如果当前消息需要你回复，只输出你要发送的聊天内容。
9. 最终输出结果中可以直接使用@用户名来称呼具体的用户，不要输出任何用户ID，也不要输出{@用户ID}这类占位符。
10. 输出必须是纯文本，不要使用Markdown、列表、标题、代码块或任何格式化。
11. 不要解释你的判断过程。
12. 不要重复历史消息，除非回复确实需要引用。
13. 回复应自然、简洁，符合私聊语境。
14. 如果任一用户明确询问你、请求你、点名@AI、或明显继续与你对话，通常需要回复。
15. 如果用户之间只是在互相聊天、自言自语、转发内容、发无关表情，或对话已经自然结束，通常不需要回复。

请根据以下历史消息判断是否需要回复，并输出最终结果：
`

func ChatAI(ctx context.Context, history string) (bool, string, error) {
	cfg := openai.DefaultConfig(config.AppConfig.AI.Key)
	cfg.BaseURL = config.AppConfig.AI.BaseURL
	client := openai.NewClientWithConfig(cfg)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: config.AppConfig.AI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: ChatPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: history,
				},
			},
		},
	)

	if err != nil {
		return false, "", err
	}

	reply := resp.Choices[0].Message.Content

	if reply == "noreply" {
		return false, "", nil
	}
	return true, reply, nil
}
