package main

import (
	"colorspacer/db"
	"colorspacer/entities"
	"colorspacer/model"
	"colorspacer/query"
	"colorspacer/types"

	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
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

	entities.Initialize(win, &clicked)                         // pixelgl had to Run() to initialize a window to initialize entities
	defaultDispatch := Dispatch{Buttons: entities.AllEntities} // AllEntities isn't initialized until entities.Initialize()

	frameTimes := []time.Time{time.Now()}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 100), basicAtlas)

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		basicTxt.Clear()

		// TODO: make fps counter an entity
		frameTimes = append(frameTimes, time.Now())
		for frameTimes[len(frameTimes)-1].Sub(frameTimes[0]).Seconds() >= 1 {
			frameTimes = frameTimes[1:]
		}

		click := win.JustPressed(pixelgl.MouseButton1)
		deltatime := frameTimes[len(frameTimes)-1].Sub(frameTimes[len(frameTimes)-2])
		defaultDispatch.Update(deltatime)
		if click {
			clicked.MousePos = win.MousePosition()
			defaultDispatch.Handle(clicked)
		}
		defaultDispatch.Draw(win) // Click indicators only work if update, then handle, then draw (or a rotation thereof)

		// Draw FPS tracker
		fmt.Fprintln(basicTxt, len(frameTimes))

		// Draw step counter
		// TODO: actually use formatting
		fmt.Fprint(basicTxt, "Step ")
		fmt.Fprintln(basicTxt, len(entities.ChosenTestColors)%7)
		basicTxt.Draw(win, pixel.IM)

		win.Update()
	}
}
