package followerDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

type FollowerDao struct {
	q *query.Query
}

func NewFollowerDao(db *gorm.DB) *FollowerDao {
	return &FollowerDao{q: query.Use(db)}
}
