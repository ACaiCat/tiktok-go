package followerdao

import (
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type FollowerDao struct {
	q *query.Query
}

func NewFollowerDao(db *gorm.DB) *FollowerDao {
	return &FollowerDao{q: query.Use(db)}
}

func (f *FollowerDao) WithTx(tx *gorm.DB) *FollowerDao {
	return &FollowerDao{q: query.Use(tx)}
}
