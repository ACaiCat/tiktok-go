package chat

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	service "github.com/ACaiCat/tiktok-go/biz/service/chat"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
	"github.com/gorilla/websocket"
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
		log.Println(err)
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

	chatService := service.NewChatService(ctx, m)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		chatService.HandleMessage(userID, string(message))
	}
}
