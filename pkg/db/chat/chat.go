package chatdao

import (
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type ChatDao struct {
	q *query.Query
}

func NewChatDao(db *gorm.DB) *ChatDao {
	return &ChatDao{q: query.Use(db)}
}

func (c *ChatDao) WithTx(tx *gorm.DB) *ChatDao {
	return &ChatDao{q: query.Use(tx)}
}
