package service

func (s *ChatService) sendMessageToUser(userID int64, msgType int, body any) (bool, error) {
	user, online := s.manager.GetOnlineUser(userID)
	if !online {
		return false, nil
	}

	if err := user.SendMessage(msgType, body); err != nil {
		return true, err
	}

	return true, nil
}
