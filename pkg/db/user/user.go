package userDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

type UserDao struct {
	q *query.Query
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{q: query.Use(db)}
}
