package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gorilla/websocket"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	service "github.com/ACaiCat/tiktok-go/biz/service/chat"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
)

var m = ws.NewOnlineUserManager()

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

func Chat(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(constants.AuthHeader)
	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := jwt.ValidateToken(token, constants.TypeAccessToken)
	if err != nil {
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	m.AddOnlineUser(userID, c)
	defer m.RemoveOnlineUser(userID)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		handleMessage(ctx, userID, string(message))
	}
}

func handleMessage(ctx context.Context, userID int64, messageText string) {
	s := service.NewChatService(ctx, m)
	var message ws.Message
	if err := json.Unmarshal([]byte(messageText), &message); err != nil {
		s.SendErr(userID, errno.ChatMsgParseErr.WithMessage(err.Error()))
		return
	}

	switch message.Type {
	case ws.MessageTypeChat:
		var chatMessage ws.ChatMessage
		if err := json.Unmarshal(message.Body, &chatMessage); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr)
			return
		}
		if err := s.HandleChatMessage(userID, &chatMessage); err != nil {
			s.SendErr(userID, errno.ServiceErr)
			hlog.CtxErrorf(ctx, "handleChatMessage: %v", err)
		}

	case ws.MessageTypeUnread:
		var unreadRequest ws.UnreadRequest
		if err := json.Unmarshal(message.Body, &unreadRequest); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr)
			return
		}
		if err := s.HandleUnreadMessage(userID, &unreadRequest); err != nil {
			s.SendErr(userID, errno.ServiceErr)
			hlog.CtxErrorf(ctx, "handleUnreadMessage: %v", err)
		}

	case ws.MessageTypeHistory:
		var historyRequest ws.HistoryRequest
		if err := json.Unmarshal(message.Body, &historyRequest); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr)
			return
		}
		if err := s.HandleHistoryMessage(userID, &historyRequest); err != nil {
			s.SendErr(userID, errno.ServiceErr)
			hlog.CtxErrorf(ctx, "handleHistoryMessage: %v", err)
		}

	default:
		s.SendErr(userID, errno.ChatMsgTypeErr)
	}
}
