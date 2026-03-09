package db

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
	"gorm.io/gorm"
)

func SetDB(db *gorm.DB) {
	query.SetDefault(db)
}
