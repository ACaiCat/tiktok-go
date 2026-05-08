package commentdao

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *CommentDao) AddVideoComment(ctx context.Context, userID int64, videoID int64, content string) error {
	var err error

	comment := model.Comment{
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}

	err = c.q.Comment.WithContext(ctx).Create(&comment)
	if err != nil {
		return errors.Wrapf(err, "AddVideoComment failed, senderID: %d, receiverID: %d", userID, videoID)
	}

	return nil
}

func (c *CommentDao) AddCommentReply(ctx context.Context, userID int64, videoID int64, commentID int64, content string) error {
	var err error

	comment := model.Comment{
		UserID:   userID,
		ParentID: new(commentID),
		Content:  content,
		VideoID:  videoID,
	}

	err = c.q.Comment.WithContext(ctx).Create(&comment)
	if err != nil {
		return errors.Wrapf(err, "AddCommentReply failed, senderID: %d, receiverID: %d", userID, videoID)
	}

	return nil
}
