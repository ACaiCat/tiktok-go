package videoDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

type VideoDao struct {
	q *query.Query
}

func NewVideoDao(db *gorm.DB) *VideoDao {
	return &VideoDao{q: query.Use(db)}
}
