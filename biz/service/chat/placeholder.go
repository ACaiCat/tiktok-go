package service

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

var userPlaceholderPattern = regexp.MustCompile(`\{@(\d+)}`)

func replaceUserPlaceholders(content string, resolveUsername func(userID int64) (string, bool)) string {
	if content == "" {
		return content
	}

	return userPlaceholderPattern.ReplaceAllStringFunc(content, func(placeholder string) string {
		matches := userPlaceholderPattern.FindStringSubmatch(placeholder)
		if len(matches) != 2 {
			return placeholder
		}

		userID, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return placeholder
		}

		username, ok := resolveUsername(userID)
		if !ok || username == "" {
			return placeholder
		}

		return "@" + username
	})
}

func (s *ChatService) replaceUserPlaceholders(content string) string {
	return replaceUserPlaceholders(content, s.newUsernameResolver())
}

func (s *ChatService) newUsernameResolver() func(userID int64) (string, bool) {
	cache := make(map[int64]string)
	missing := make(map[int64]struct{})

	return func(userID int64) (string, bool) {
		if username, ok := cache[userID]; ok {
			return username, true
		}
		if _, ok := missing[userID]; ok {
			return "", false
		}

		user, err := s.userDao.GetByID(s.ctx, userID)
		if err != nil {
			log.Println("failed to get user for placeholder:", err)
			missing[userID] = struct{}{}
			return "", false
		}
		if user == nil || strings.TrimSpace(user.Username) == "" {
			missing[userID] = struct{}{}
			return "", false
		}

		cache[userID] = user.Username
		return user.Username, true
	}
}

func normalizeAIContent(content string, resolveUsername func(userID int64) (string, bool)) string {
	if content == "" {
		return content
	}

	content = strings.ReplaceAll(content, "{@AI}", "@AI")
	return replaceUserPlaceholders(content, resolveUsername)
}

func buildAIHistory(messages []*model.ChatMessage, resolveUsername func(userID int64) (string, bool)) string {
	var history strings.Builder

	for _, message := range messages {
		identity := "@AI"
		if !message.IsAi {
			if username, ok := resolveUsername(message.SenderID); ok && username != "" {
				identity = "@" + username
			} else {
				identity = "@" + strconv.FormatInt(message.SenderID, 10)
			}
		}

		history.WriteString(identity)
		history.WriteString(": ")
		history.WriteString(normalizeAIContent(message.Content, resolveUsername))
		history.WriteByte('\n')
	}

	return history.String()
}

func (s *ChatService) buildAIHistory(messages []*model.ChatMessage) string {
	return buildAIHistory(messages, s.newUsernameResolver())
}
