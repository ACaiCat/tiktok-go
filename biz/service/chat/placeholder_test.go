package service

import (
	"testing"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestReplaceUserPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		resolver func(int64) (string, bool)
		want     string
	}{
		{
			name:    "single placeholder",
			content: "你好 {@1}",
			resolver: func(userID int64) (string, bool) {
				if userID == 1 {
					return "alice", true
				}
				return "", false
			},
			want: "你好 @alice",
		},
		{
			name:    "multiple and repeated placeholders",
			content: "{@1} 找 {@2}，再找 {@1}",
			resolver: func(userID int64) (string, bool) {
				switch userID {
				case 1:
					return "alice", true
				case 2:
					return "bob", true
				default:
					return "", false
				}
			},
			want: "@alice 找 @bob，再找 @alice",
		},
		{
			name:    "unknown user keeps placeholder",
			content: "你好 {@404}",
			resolver: func(int64) (string, bool) {
				return "", false
			},
			want: "你好 {@404}",
		},
		{
			name:    "invalid placeholder format stays unchanged",
			content: "你好 {@abc} 和 {@AI}",
			resolver: func(int64) (string, bool) {
				return "nobody", true
			},
			want: "你好 {@abc} 和 {@AI}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceUserPlaceholders(tt.content, tt.resolver)
			if got != tt.want {
				t.Fatalf("replaceUserPlaceholders() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildAIHistory(t *testing.T) {
	messages := []*model.ChatMessage{
		{
			SenderID: 1,
			Content:  "你好 {@2}",
			IsAi:     false,
		},
		{
			SenderID: 99,
			Content:  "{@1} 你好，我是 {@AI}",
			IsAi:     true,
		},
		{
			SenderID: 2,
			Content:  "收到",
			IsAi:     false,
		},
	}

	resolver := func(userID int64) (string, bool) {
		switch userID {
		case 1:
			return "alice", true
		case 2:
			return "bob", true
		default:
			return "", false
		}
	}

	want := "@alice: 你好 @bob\n@AI: @alice 你好，我是 @AI\n@bob: 收到\n"
	got := buildAIHistory(messages, resolver)
	if got != want {
		t.Fatalf("buildAIHistory() = %q, want %q", got, want)
	}
}

func TestBuildAIHistoryFallbackToUserID(t *testing.T) {
	messages := []*model.ChatMessage{{
		SenderID: 404,
		Content:  "你好",
		IsAi:     false,
	}}

	got := buildAIHistory(messages, func(int64) (string, bool) {
		return "", false
	})
	want := "@404: 你好\n"
	if got != want {
		t.Fatalf("buildAIHistory() fallback = %q, want %q", got, want)
	}
}
