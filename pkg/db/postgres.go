package db

import (
	"fmt"

	"github.com/ACaiCat/tiktok-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.AppConfig.Postgres.Host,
		config.AppConfig.Postgres.User,
		config.AppConfig.Postgres.Password,
		config.AppConfig.Postgres.DBName,
		config.AppConfig.Postgres.Port,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
