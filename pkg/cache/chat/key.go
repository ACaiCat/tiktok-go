package chatcache

import "fmt"

func normalizeConversationUserIDs(userID int64, otherUserID int64) (int64, int64) {
	if userID > otherUserID {
		return otherUserID, userID
	}
	return userID, otherUserID
}

func getHistoryKey(userID int64, otherUserID int64, pageSize int, pageNum int) string {
	left, right := normalizeConversationUserIDs(userID, otherUserID)
	return fmt.Sprintf("chat:history:%d:%d:%d:%d", left, right, pageSize, pageNum)
}

func getHistoryKeyPattern(userID int64, otherUserID int64) string {
	left, right := normalizeConversationUserIDs(userID, otherUserID)
	return fmt.Sprintf("chat:history:%d:%d:*", left, right)
}

func getUnreadKey(userID int64, senderID int64) string {
	return fmt.Sprintf("chat:unread:%d:%d", userID, senderID)
}
