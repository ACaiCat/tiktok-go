package chatDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

type ChatDao struct {
	q *query.Query
}

func NewChatDao(db *gorm.DB) *ChatDao {
	return &ChatDao{q: query.Use(db)}
}
