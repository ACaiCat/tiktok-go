package likedao

import (
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type LikeDao struct {
	q *query.Query
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{q: query.Use(db)}
}
