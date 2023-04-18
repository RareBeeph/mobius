package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

// type coloredRect struct {
// 	bounds pixel.Rect
// 	color  pixel.RGBA
// }

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

	button := imdraw.New(nil)

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

	// TODO: Draw those rectangles to the window

	var frameTimes []time.Time

	// Moved up here so that we're not re-initializing their memory on every frame
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 100), basicAtlas)

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		frameTimes = append(frameTimes, time.Now())
		for frameTimes[len(frameTimes)-1].Sub(frameTimes[0]).Seconds() >= 1 {
			frameTimes = frameTimes[1:]
		}

		fmt.Fprintln(basicTxt, len(frameTimes))
		basicTxt.Draw(win, pixel.IM)

		// TODO: Update click logic to use receiver method
		if win.JustPressed(pixelgl.MouseButton1) {
			mpos := win.MousePosition()
			for i, r := range rects {
				if mpos.X >= r.Bounds.Min.X && mpos.X <= r.Bounds.Max.X && mpos.Y >= r.Bounds.Min.Y && mpos.Y <= r.Bounds.Max.Y {
					butt := imdraw.New(nil)
					pushRectToImd(pixel.V((float64)(i*10), 0), pixel.V((float64)(i*10+10), 10), r.Color, butt)
					butt.Draw(win)
				}
			}
			butt2 := imdraw.New(nil)
			pushRectToImd(pixel.V(0, 10), pixel.V(10, 20), pixel.RGB(0, 1, 0), butt2)
			butt2.Draw(win)
		}

		button.Draw(win)
		win.Update()
	}
}

func pushRectToImd(p0 pixel.Vec, p1 pixel.Vec, color pixel.RGBA, imd *imdraw.IMDraw) {
	imd.Color = color
	imd.Push(p0)
	imd.Push(p1)
	imd.Rectangle(0)
}
