package videodao

import (
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type VideoDao struct {
	q *query.Query
}

func NewVideoDao(db *gorm.DB) *VideoDao {
	return &VideoDao{q: query.Use(db)}
}
