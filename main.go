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

	var clicked = types.Event{EventType: types.Click}
	var delta = types.Event{EventType: types.Drag}

	entities.Initialize(win, &clicked) // pixelgl had to Run() to initialize a window to initialize entities

	defaultDispatch := Dispatch{Entities: []types.E{entities.Scene}}

	thisPos := win.MousePosition()
	var lastPos pixel.Vec

	var lastFrame time.Time
	thisFrame := time.Now()
	var deltatime time.Duration

	ButtonsUsed := []pixelgl.Button{pixelgl.MouseButton1, pixelgl.MouseButton2, pixelgl.KeyC}

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		lastFrame = thisFrame
		thisFrame = time.Now()
		deltatime = thisFrame.Sub(lastFrame)

		lastPos = thisPos
		thisPos = win.MousePosition()
		delta.MouseVel = thisPos.Sub(lastPos)
		delta.MousePos = thisPos

		clicked.Buttons = []pixelgl.Button{}
		delta.Buttons = []pixelgl.Button{}
		for _, b := range ButtonsUsed {
			if win.JustPressed(b) {
				clicked.Buttons = append(clicked.Buttons, b)
				delta.InitialPos = win.MousePosition()
			}
			if win.Pressed(b) {
				delta.Buttons = append(delta.Buttons, b)
			}
		}

		defaultDispatch.Update(deltatime)
		clicked.MousePos = win.MousePosition()
		defaultDispatch.Handle(&clicked)
		defaultDispatch.Handle(&delta) // Could be fused into clicked if events stored separate slices for clicked and held buttons

		defaultDispatch.Draw(win) // Click indicators only work if update, then handle, then draw (or a rotation thereof)

		win.Update()
	}
}
