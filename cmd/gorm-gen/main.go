package main

import (
	"gorm.io/gen"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/db"
)

func main() {
	config.Init()

	g := gen.NewGenerator(gen.Config{
		OutPath:       "pkg/db/query",
		Mode:          gen.WithoutContext | gen.WithQueryInterface,
		FieldNullable: true,
	})

	db.InitPostgres()
	g.UseDB(db.DB)

	videoModel := g.GenerateModel("videos")

	commentModel := g.GenerateModel("comments")
	userModel := g.GenerateModel("users")
	likeModel := g.GenerateModel("likes")
	followerModel := g.GenerateModel("followers")
	messageModel := g.GenerateModel("chat_messages")

	g.ApplyBasic(videoModel, userModel, likeModel, commentModel, followerModel, messageModel)

	g.Execute()
}
