package commentDao

import (
	"errors"
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"gorm.io/gorm"
)

func (c *CommentDao) GetCommentCount(videoID int64) (int64, error) {
	var err error

	count, err := c.q.Comment.
		Where(c.q.Comment.VideoID.Eq(videoID)).
		Count()
	if err != nil {
		log.Printf("failed to get comment count for videoID %d: %v", videoID, err)
		return 0, err
	}

	return count, nil
}

func (c *CommentDao) GetCommentCounts(videoIDs []int64) (map[int64]int64, error) {
	var err error

	type Result struct {
		VideoID int64 `gorm:"column:video_id"`
		Count   int64 `gorm:"column:count"`
	}

	var results []Result

	err = c.q.Comment.
		Select(c.q.Comment.VideoID, c.q.Comment.ID.Count().As("count")).
		Where(c.q.Comment.VideoID.In(videoIDs...)).
		Group(c.q.Comment.VideoID).
		Scan(&results)

	if err != nil {
		log.Printf("failed to get comment counts for videoIDs %v: %v", videoIDs, err)
		return nil, err
	}
	commentMap := make(map[int64]int64)
	for _, r := range results {
		commentMap[r.VideoID] = r.Count
	}

	return commentMap, nil
}

func (c *CommentDao) GetCommentByID(commentID int64) (*model.Comment, error) {
	var err error

	subComment := c.q.Comment.As("sub_comment")

	comment, err := c.q.Comment.
		Select(c.q.Comment.ALL,
			c.q.Like.ID.Count().As("like_count"),
			c.q.Comment.ID.Count().As("child_count"),
		).
		LeftJoin(c.q.Like, c.q.Like.CommentID.EqCol(c.q.Comment.ID)).
		LeftJoin(subComment, subComment.ParentID.EqCol(c.q.Comment.ID)).
		Group(c.q.Comment.ID).
		Where(c.q.Comment.ID.Eq(commentID)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Printf("failed to get comment by ID %d: %v", commentID, err)
		return nil, err
	}

	return comment, nil

}

func (c *CommentDao) GetCommentsByVideoID(videoID int64, pageSize int, pageNum int) ([]*model.Comment, error) {
	var err error

	subComment := c.q.Comment.As("sub_comment")

	comments, err := c.q.Comment.
		Select(
			c.q.Comment.ALL,
			c.q.Like.ID.Count().As("like_count"),
			c.q.Comment.ID.Count().As("child_count"),
		).
		LeftJoin(c.q.Like, c.q.Like.CommentID.EqCol(c.q.Comment.ID)).
		LeftJoin(subComment, subComment.ParentID.EqCol(c.q.Comment.ID)).
		Group(c.q.Comment.ID).
		Where(c.q.Comment.VideoID.Eq(videoID)).
		Order(c.q.Comment.CreatedAt.Desc()).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Printf("failed to get comments by video ID %d: %v", videoID, err)
		return nil, err
	}

	return comments, nil

}
