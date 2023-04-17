package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type coloredRect struct {
	bounds pixel.Rect
	color  pixel.RGBA
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

	button := imdraw.New(nil)

	var rects []coloredRect
	var r coloredRect

	r.bounds = pixel.R(200, 100, 500, 300)
	r.color = pixel.RGB(1, 0, 0)
	rects = append(rects, r)
	pushRectToImd(rects[0].bounds.Min, rects[0].bounds.Max, rects[0].color, button)

	var test []time.Time

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		append(test, time.Now())

		if win.JustPressed(pixelgl.MouseButton1) {
			mpos := win.MousePosition()
			for _, r := range rects {
				if mpos.X >= r.bounds.Min.X && mpos.X <= r.bounds.Max.X && mpos.Y >= r.bounds.Min.Y && mpos.Y <= r.bounds.Max.Y {
					butt := imdraw.New(nil)
					pushRectToImd(pixel.V(0, 0), pixel.V(10, 10), pixel.RGB(1, 0, 0), butt)
					butt.Draw(win)
				}
			}
			butt2 := imdraw.New(nil)
			pushRectToImd(pixel.V(10, 0), pixel.V(20, 10), pixel.RGB(0, 1, 0), butt2)
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
