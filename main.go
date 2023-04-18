package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

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

	var rects []ColoredRect = []ColoredRect{
		{
			Bounds: pixel.R(200, 100, 500, 300),
			Color:  pixel.RGB(1, 0, 0),
		},

		{
			Bounds: pixel.R(500, 100, 600, 200),
			Color:  pixel.RGB(0.7, 0.4, 0.2),
		},
	}

	var frameTimes []time.Time

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 100), basicAtlas)

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		basicTxt.Clear()

		frameTimes = append(frameTimes, time.Now())
		for frameTimes[len(frameTimes)-1].Sub(frameTimes[0]).Seconds() >= 1 {
			frameTimes = frameTimes[1:]
		}

		fmt.Fprintln(basicTxt, len(frameTimes))
		basicTxt.Draw(win, pixel.IM)

		mb1 := win.JustPressed(pixelgl.MouseButton1)
		for i, r := range rects {
			r.Draw(win)
			if mb1 && r.Contains(win.MousePosition()) {
				(&ColoredRect{Bounds: pixel.R((float64)(i*10), 0, (float64(i*10) + 10), 10), Color: r.Color}).Draw(win)
			}
		}
		if mb1 {
			(&ColoredRect{Bounds: pixel.R(0, 10, 10, 20), Color: pixel.RGB(0, 1, 0)}).Draw(win)
		}

		win.Update()
	}
}
