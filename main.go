package main

import (
	"fmt"
	"math/rand"
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

	controlRects := []ColoredRect{
		{
			Bounds: pixel.R(200, 100, 500, 300),
			Color:  chooseControlColor(),
		},
		{
			Bounds: pixel.R(600, 100, 900, 300),
			Color:  chooseControlColor(),
		},
	}

	/*var testRects []ColoredRect = []ColoredRect{
		{
			Bounds: pixel.R(500, 100, 600, 200),
			Color:  pixel.RGB(0.7, 0.4, 0.2),
		},
	}*/
	depth := 0
	step := 0
	testRects := makeTestRects(controlRects, depth, step, []pixel.RGBA{})

	var frameTimes []time.Time

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 100), basicAtlas)

	clickIndicator := &ColoredRect{Bounds: pixel.R(0, 10, 10, 20), Color: pixel.RGB(0, 1, 0)}
	collisionIndicator := &ColoredRect{Bounds: pixel.R(0, 0, 10, 10), Color: pixel.RGB(0, 1, 0)}

	chosenTestColors := []pixel.RGBA{}

	for !win.Closed() {
		// Blanking
		win.Clear(pixel.RGB(0, 0, 0))
		basicTxt.Clear()

		// Logic sweep
		frameTimes = append(frameTimes, time.Now())
		for frameTimes[len(frameTimes)-1].Sub(frameTimes[0]).Seconds() >= 1 {
			frameTimes = frameTimes[1:]
		}

		mouseClicked := win.JustPressed(pixelgl.MouseButton1)

		for idx, rect := range append(controlRects, testRects...) {
			rect.Draw(win)
			if mouseClicked && rect.Contains(win.MousePosition()) {
				// Indicate which color was clicked
				collisionIndicator.Bounds = pixel.R(float64(idx*10), 0, float64(idx*10)+10, 10)
				collisionIndicator.Color = rect.Color
				collisionIndicator.Draw(win)
				chosenTestColors = append(chosenTestColors, rect.Color)

				// Generate a new set of colors to compare
				// controlRects[0].Color = chooseControlColor()
				step++
				testRects = makeTestRects(controlRects, depth, step, chosenTestColors)
			}
		}
		if mouseClicked {
			clickIndicator.Draw(win)
		}

		// Draw FPS tracker
		fmt.Fprintln(basicTxt, len(frameTimes))
		basicTxt.Draw(win, pixel.IM)

		// Send to screen!
		win.Update()
	}
}

func chooseControlColor() pixel.RGBA {
	return pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
}

/*func chooseTestColors(ctrlColor pixel.RGBA) []pixel.RGBA {
	out := []pixel.RGBA{}
	idx := 0
	ctrlColor.R = pixel.Clamp(ctrlColor.R, 0.2, 0.8)
	ctrlColor.G = pixel.Clamp(ctrlColor.G, 0.2, 0.8)
	ctrlColor.B = pixel.Clamp(ctrlColor.B, 0.2, 0.8)
	for idx < 5 {
		out = append(out, pixel.RGB(ctrlColor.R+rand.Float64()*0.2, ctrlColor.G+rand.Float64()*0.2, ctrlColor.B+rand.Float64()*0.2))
		idx++
	}
	return out
}*/

func firstChooseTestColors(inColor pixel.RGBA, depth int, step int, ctc []pixel.RGBA) []pixel.RGBA {
	out := []pixel.RGBA{}
	if step < 4 {
		offset := 1 / (float64(int(4) << depth))
		out = append(out, pixel.RGB(inColor.R-offset, inColor.G-offset, inColor.B-offset))

		if step&1 == 1 {
			out[0].R += 2 * offset
		}
		if step&2 == 2 {
			out[0].G += 2 * offset
		}
		out = append(out, out[0])
		out[1].B += 2 * offset
	} else if step == 4 {
		out = append(out, ctc[0], ctc[1])
	} else if step == 5 {
		out = append(out, ctc[2], ctc[3])
	} else if step == 6 {
		out = append(out, ctc[4], ctc[5])
	} else {
		out = firstChooseTestColors(ctc[6], depth+1, step-7, ctc[6:])
	}

	return out
}

func makeTestRects(controlRects []ColoredRect, depth int, step int, ctc []pixel.RGBA) []ColoredRect {
	var testRects []ColoredRect
	for idx, col := range firstChooseTestColors(pixel.RGB(0.5, 0.5, 0.5), depth, step, ctc) {
		rect := ColoredRect{Bounds: pixel.R(500, 100+float64(idx*100), 600, 200+float64(idx*100)), Color: col}
		testRects = append(testRects, rect)
	}
	return testRects
}
