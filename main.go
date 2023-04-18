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

	r.bounds = pixel.R(500, 100, 600, 200)
	r.color = pixel.RGB(0.7, 0.4, 0.2)
	rects = append(rects, r)

	for _, r := range rects {
		pushRectToImd(r.bounds.Min, r.bounds.Max, r.color, button)
	}
	var test []time.Time

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		test = append(test, time.Now())
		for test[len(test)-1].Sub(test[0]).Seconds() >= 1 {
			test = test[1:]
		}
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 100), basicAtlas)

		fmt.Fprintln(basicTxt, len(test))
		basicTxt.Draw(win, pixel.IM)

		if win.JustPressed(pixelgl.MouseButton1) {
			mpos := win.MousePosition()
			for i, r := range rects {
				if mpos.X >= r.bounds.Min.X && mpos.X <= r.bounds.Max.X && mpos.Y >= r.bounds.Min.Y && mpos.Y <= r.bounds.Max.Y {
					butt := imdraw.New(nil)
					pushRectToImd(pixel.V((float64)(i*10), 0), pixel.V((float64)(i*10+10), 10), r.color, butt)
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
