package service

import "github.com/ACaiCat/tiktok-go/pkg/errno"

func (s *ChatService) sendMessage(userID int64, msgType int, body any) (bool, error) {
	user, online := s.manager.GetOnlineUser(userID)
	if !online {
		return false, nil
	}

	if err := user.SendMessage(msgType, body); err != nil {
		return true, err
	}

	return true, nil
}

func (s *ChatService) SendErr(userID int64, err errno.ErrNo) {
	if u, online := s.manager.GetOnlineUser(userID); online {
		u.SendError(int(err.ErrCode), err.ErrMsg)
	}
}
