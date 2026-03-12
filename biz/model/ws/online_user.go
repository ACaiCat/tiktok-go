package ws

import (
	"github.com/gorilla/websocket"
)

type OnlineUser struct {
	UserID int64
	Conn   *websocket.Conn
}

func (u *OnlineUser) SendMessage(msgType int, messageBody any) error {
	message := SendMessage{
		Type: msgType,
		Body: messageBody,
	}

	return u.Conn.WriteJSON(&message)
}

func (u *OnlineUser) SendError(code int, msg string) {
	_ = u.SendMessage(MessageTypeError, &ErrorMessage{Code: code, Message: msg})
}
