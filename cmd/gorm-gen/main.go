package main

import (
	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func main() {
	config.Init()

	g := gen.NewGenerator(gen.Config{
		OutPath: "pkg/db/query",
		Mode:    gen.WithoutContext | gen.WithQueryInterface,
	})

	db.InitPostgres()
	g.UseDB(db.DB)

	videoModel := g.GenerateModel("videos",
		gen.FieldNew("CommentCount", "int64", field.Tag{
			"gorm": "column:comment_count;->",
			"json": "comment_count",
		}),
		gen.FieldNew("LikeCount", "int64", field.Tag{
			"gorm": "column:like_count;->",
			"json": "like_count",
		}),
	)

	commentModel := g.GenerateModel("comments",
		gen.FieldNew("ChildCount", "int64", field.Tag{
			"gorm": "column:child_count;->",
			"json": "child_count",
		}),
		gen.FieldNew("LikeCount", "int64", field.Tag{
			"gorm": "column:like_count;->",
			"json": "like_count",
		}),
	)
	userModel := g.GenerateModel("users")
	likeModel := g.GenerateModel("likes")

	followerModel := g.GenerateModel("followers")

	g.ApplyBasic(videoModel, userModel, likeModel, commentModel, followerModel)

	g.Execute()
}
