package main

import (
	"colorspacer/entities"
	"colorspacer/types"

	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "test",
		Bounds: pixel.R(0, 0, 500, 450),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var clicked = types.Event{EventType: types.Click}
	var delta = types.Event{EventType: types.Drag}

	entities.InitSceneTwo(win, &clicked) // pixelgl had to Run() to initialize a window to initialize entities

	defaultDispatch := Dispatch{Entities: []types.EI{entities.Scene2}}

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
