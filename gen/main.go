package main

import (
	"gorm.io/gen"

	"colorspacer/db"
	"colorspacer/db/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../db/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db.Connection)

	// Generate basic type-safe DAO API for model structs following conventions
	g.ApplyBasic(model.AllModels...)

	// Generate the code
	g.Execute()
}
