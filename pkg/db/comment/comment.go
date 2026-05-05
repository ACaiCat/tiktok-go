package commentdao

import (
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type CommentDao struct {
	q *query.Query
}

func NewCommentDao(db *gorm.DB) *CommentDao {
	return &CommentDao{q: query.Use(db)}
}

func (c *CommentDao) WithTx(tx *gorm.DB) *CommentDao {
	return &CommentDao{q: query.Use(tx)}
}
