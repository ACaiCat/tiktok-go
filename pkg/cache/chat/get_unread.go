package chatcache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatCache) GetUnreadMessages(ctx context.Context, userID int64, senderID int64) ([]*model.ChatMessage, error) {
	data, err := c.c.Get(ctx, getUnreadKey(userID, senderID)).Bytes()
	if err != nil {
		return nil, err
	}

	var messages []*model.ChatMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, errors.Wrapf(err, "GetUnreadMessages failed,  bucket=%s, object=%s", constants.AvatarBucketName, data)
	}

	return messages, nil
}
