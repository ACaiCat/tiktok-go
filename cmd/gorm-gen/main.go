package main

import (
	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	"gorm.io/gen"
)

func main() {
	config.Init()

	g := gen.NewGenerator(gen.Config{
		OutPath: "pkg/db/query",
		Mode:    gen.WithoutContext | gen.WithQueryInterface,
	})

	db.InitPostgres()
	g.UseDB(db.DB)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
