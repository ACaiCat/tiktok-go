package chat

import (
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
	m.AddOnlineUser(userID, c)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			m.RemoveOnlineUser(userID)
			break
		}
		service.NewChatService(m).HandleMessage(userID, string(message))
	}
}
