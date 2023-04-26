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
	db.Connection.AutoMigrate(model.AllModels...)
}

func main() {
	pixelgl.Run(run)
}

var gstep int
var gdepth int
var chosenTestColors []pixel.RGBA
var testRects []*Button

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

	m := query.Midpoint

	rand.Seed(time.Now().UnixMicro())

	controlRects := []Button{
		{ColoredRect: ColoredRect{
			Bounds: pixel.R(200, 100, 500, 300),
			Color:  chooseControlColor(),
		}, procedure: func() {}},
		{ColoredRect: ColoredRect{
			Bounds: pixel.R(600, 100, 900, 300),
			Color:  chooseControlColor(),
		}, procedure: func() {}},
	}

	gdepth = 0
	gstep = 0
	testRects = makeTestRects(controlRects, gdepth, gstep, []pixel.RGBA{})

	saveButton := Button{
		ColoredRect: ColoredRect{
			Bounds: pixel.R(400, 400, 700, 600),
			Color:  pixel.RGB(0.8, 0.8, 0.8),
		},
		procedure: func() {
			if len(chosenTestColors) > 0 {
				start := model.NewColorFromRgba(controlRects[0].Color)
				end := model.NewColorFromRgba(controlRects[1].Color)
				mid := model.NewColorFromRgba(chosenTestColors[len(chosenTestColors)-1])
				a := model.Midpoint{
					Startpoint: *start,
					Endpoint:   *end,
					Midpoint:   *mid,
				}
				m.Create(&a)

				//Debug
				b, _ := m.Preload(m.Midpoint).Last()
				log.Printf("ID: %d, R: %f, G: %f, B: %f", b.Midpoint.ID, b.Midpoint.R, b.Midpoint.G, b.Midpoint.B)
			}
		},
	}

	var frameTimes []time.Time

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 100), basicAtlas)

	// TODO: make these work again
	clickIndicator := Button{ColoredRect: ColoredRect{Bounds: pixel.R(0, 10, 10, 20), Color: pixel.RGB(0, 1, 0)}, procedure: func() {}}
	collisionIndicator := Button{ColoredRect: ColoredRect{Bounds: pixel.R(0, 0, 10, 10), Color: pixel.RGB(0, 1, 0)}, procedure: func() {}}

	entities := []*Button{}
	entities = append(entities, &controlRects[0], &controlRects[1])
	entities = append(entities, testRects...)
	entities = append(entities, &saveButton, &clickIndicator, &collisionIndicator)

	clicked := Event{}

	for !win.Closed() {
		/*
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
				start := model.NewColorFromRgba(controlRects[0].Color)
				end := model.NewColorFromRgba(controlRects[1].Color)
				mid := model.NewColorFromRgba(chosenTestColors[len(chosenTestColors)-1])
				a := model.Midpoint{
					Startpoint: *start,
					Endpoint:   *end,
					Midpoint:   *mid,
				}
				m.Create(&a)

				//Debug
				b, _ := m.Preload(m.Midpoint).Last()
				log.Printf("ID: %d, R: %f, G: %f, B: %f", b.Midpoint.ID, b.Midpoint.R, b.Midpoint.G, b.Midpoint.B)
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
		*/

		win.Clear(pixel.RGB(0, 0, 0))
		basicTxt.Clear()

		click := win.JustPressed(pixelgl.MouseButton1)
		if click {
			clicked.mousePos = win.MousePosition()
		}
		for _, e := range entities {
			e.Draw(win)
			if click {
				e.Handle(clicked)
			}
		}

		frameTimes = append(frameTimes, time.Now())
		for frameTimes[len(frameTimes)-1].Sub(frameTimes[0]).Seconds() >= 1 {
			frameTimes = frameTimes[1:]
		}

		// Draw FPS tracker
		fmt.Fprintln(basicTxt, len(frameTimes))

		// Draw step counter
		// TODO: actually use formatting
		fmt.Fprint(basicTxt, "Step ")
		fmt.Fprintln(basicTxt, gstep%7)
		basicTxt.Draw(win, pixel.IM)

		win.Update()

		/*
			The run loop handles the screen.
				win.Clear(pixel.RGB(0,0,0))
				...
				win.Update()


			Entities should be stored as a collection that can be iterated through.
				// Should I have two collections, one for interactables and one for things that just do their own thing?
				// Should the creation and addition of entities to this collection be delegated too?

				var entities []Entity
				var event Event{
					Click: win.JustPressed(pixelgl.MouseButton1)
					Mpos: win.MousePosition()
				}
				...
				for _,e := range entities {
					e.Handle(event) // Each entity decides if it cares about each aspect of the event
				}
		*/
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
		offset := 1 / (float64(int(3) << depth))
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

func makeTestRects(controlRects []Button, depth int, step int, ctc []pixel.RGBA) (tr []*Button) {
	for idx, col := range firstChooseTestColors(pixel.RGB(0.5, 0.5, 0.5), depth, step, ctc) {
		rect := Button{ColoredRect: ColoredRect{Bounds: pixel.R(500, 100+float64(idx*100), 600, 200+float64(idx*100)), Color: col}}
		rect.procedure = func() {
			// Indicate which color was clicked
			// collisionIndicator.Bounds = pixel.R(float64(idx*10), 0, float64(idx*10)+10, 10)
			// collisionIndicator.Color = rect.Color
			// collisionIndicator.Draw(win)
			chosenTestColors = append(chosenTestColors, rect.Color)

			// Generate a new set of colors to compare
			// controlRects[0].Color = chooseControlColor()
			gstep++
			a := makeTestRects(controlRects, gdepth, gstep, chosenTestColors)
			*testRects[0] = *a[0]
			*testRects[1] = *a[1]
		}
		tr = append(tr, &rect)
	}
	return tr
}
