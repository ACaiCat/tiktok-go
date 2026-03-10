package commentDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

type CommentDao struct {
	q *query.Query
}

func NewCommentDao(db *gorm.DB) *CommentDao {
	return &CommentDao{q: query.Use(db)}
}
