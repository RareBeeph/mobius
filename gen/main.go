package main

import (
	"gorm.io/gen"

	"colorspacer/db"
	"colorspacer/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db.Connection)

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.Midpoint{})
	db.Connection.AutoMigrate(model.Midpoint{})

	// Generate the code
	g.Execute()
}
