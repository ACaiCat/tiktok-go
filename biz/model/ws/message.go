package ws

import "encoding/json"

const (
	MessageTypeChat    = 1
	MessageTypeHistory = 2
	MessageTypeUnread  = 3
	MessageTypeError   = 4
)

type Message struct {
	Type int             `json:"type"`
	Body json.RawMessage `json:"body"`
}

type SendMessage struct {
	Type int `json:"type"`
	Body any `json:"body"`
}

type ChatMessage struct {
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	IsAI       bool   `json:"is_ai"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
}

type HistoryMessage struct {
	Messages []*ChatMessage `json:"messages"`
}

type UnreadMessage struct {
	Messages []*ChatMessage `json:"messages"`
}

type HistoryRequest struct {
	Sender   int64 `json:"user_id"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

type UnreadRequest struct {
	Sender int64 `json:"user_id"`
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
