package likeDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

type LikeDao struct {
	q *query.Query
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{q: query.Use(db)}
}
