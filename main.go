package main

import (
	"colorspacer/db"
	"colorspacer/db/model"
	"colorspacer/db/query"
	"colorspacer/entities"
	"colorspacer/types"

	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func init() {
	query.SetDefault(db.Connection)
	db.Connection.AutoMigrate(model.AllModels...)

	rand.Seed(time.Now().UnixMicro())
}

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "test",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var clicked types.Event
	entities.Initialize(win, &clicked) // pixelgl had to Run() to initialize a window to initialize entities
	defaultDispatch := Dispatch{
		Buttons:    entities.AllEntities, // AllEntities isn't initialized until entities.Initialize()
		TextFields: entities.AllTexts,
	}

	var lastFrame time.Time
	thisFrame := time.Now()
	var deltatime time.Duration

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		lastFrame = thisFrame
		thisFrame = time.Now()
		deltatime = thisFrame.Sub(lastFrame)

		click := win.JustPressed(pixelgl.MouseButton1)
		defaultDispatch.Update(deltatime)
		if click {
			clicked.MousePos = win.MousePosition()
			defaultDispatch.Handle(clicked)
		}
		defaultDispatch.Draw(win) // Click indicators only work if update, then handle, then draw (or a rotation thereof)

		win.Update()
	}
}
