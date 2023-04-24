package main

import (
	"colorspacer/db"
	"colorspacer/model"
	"colorspacer/query"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func init() {
	query.SetDefault(db.Connection)
	db.Connection.AutoMigrate(model.Midpoint{})
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

	/* q := query.Use(db.Connection)
	db.Connection.AutoMigrate(model.Midpoint{}) */

	rand.Seed(time.Now().UnixMicro())

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

	depth := 0
	step := 0
	testRects := makeTestRects(controlRects, depth, step, []pixel.RGBA{})

	saveButton := ColoredRect{
		Bounds: pixel.R(400, 400, 700, 600),
		Color:  pixel.RGB(0.8, 0.8, 0.8),
	}

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

		// Save button
		saveButton.Draw(win)
		if mouseClicked && saveButton.Contains(win.MousePosition()) && len(chosenTestColors) > 0 {
			a := model.Midpoint{
				StartpointR: controlRects[0].Color.R,
				StartpointG: controlRects[0].Color.G,
				StartpointB: controlRects[0].Color.B,
				EndpointR:   controlRects[1].Color.R,
				EndpointG:   controlRects[1].Color.G,
				EndpointB:   controlRects[1].Color.B,
				MidpointR:   chosenTestColors[len(chosenTestColors)-1].R,
				MidpointG:   chosenTestColors[len(chosenTestColors)-1].G,
				MidpointB:   chosenTestColors[len(chosenTestColors)-1].B,
			}
			m := query.Midpoint
			err := m.Create(&a) // segfaults
			log.Println(err)

			//Debug
			b, _ := m.Last()
			log.Printf("R: %f, G: %f, B: %f", b.MidpointR, b.MidpointG, b.MidpointB)
		}

		// Draw FPS tracker
		fmt.Fprintln(basicTxt, len(frameTimes))

		// Draw step counter
		// TODO: actually use formatting
		fmt.Fprint(basicTxt, "Step ")
		fmt.Fprintln(basicTxt, step%7)
		basicTxt.Draw(win, pixel.IM)

		// Send to screen!
		win.Update()
	}
}

func chooseControlColor() pixel.RGBA {
	return pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
}

func firstChooseTestColors(inColor pixel.RGBA, depth int, step int, ctc []pixel.RGBA) (out []pixel.RGBA) {
	// When this is initially called by makeTestRects, it will be centered on 0.5 gray.

	if step < 4 {
		// Over the course of 4 steps, this will generate all 8 colors at a given distance from the input color.
		// The offset will start out as 0.25 and cut in half every depth increment.
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
		// Semifinals round 1
		out = append(out, ctc[0], ctc[1])
	} else if step == 5 {
		// Semifinals round 2
		out = append(out, ctc[2], ctc[3])
	} else if step == 6 {
		// Finals
		out = append(out, ctc[4], ctc[5])
	} else {
		// Recurse with incremented depth, centered on the winner of the tournament.
		out = firstChooseTestColors(ctc[6], depth+1, step-7, ctc[7:])
	}

	return out
}

func makeTestRects(controlRects []ColoredRect, depth int, step int, ctc []pixel.RGBA) (testRects []ColoredRect) {
	for idx, col := range firstChooseTestColors(pixel.RGB(0.5, 0.5, 0.5), depth, step, ctc) {
		rect := ColoredRect{Bounds: pixel.R(500, 100+float64(idx*100), 600, 200+float64(idx*100)), Color: col}
		testRects = append(testRects, rect)
	}
	return testRects
}
