package ai

type ToolCallContext struct {
	AllowedUserIDs map[int64]struct{}
}

func NewToolCallContext(userIDs ...int64) ToolCallContext {
	allowedUserIDs := make(map[int64]struct{}, len(userIDs))
	for _, userID := range userIDs {
		allowedUserIDs[userID] = struct{}{}
	}

	return ToolCallContext{
		AllowedUserIDs: allowedUserIDs,
	}
}

func (c ToolCallContext) CanAccessUser(userID int64) bool {
	_, ok := c.AllowedUserIDs[userID]
	return ok
}
