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
	"github.com/faiface/pixel/imdraw"
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

var chosenTestColors []pixel.RGBA

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

	// TODO: make this less bloated
	controlRects := []Button{
		{ColoredRect: ColoredRect{
			Bounds: pixel.R(200, 100, 500, 300),
			Color:  chooseControlColor(),
			Entity: Entity{UpdateFunc: func(time.Duration) {}},
		}, OnEvent: func() {}},
		{ColoredRect: ColoredRect{
			Bounds: pixel.R(600, 100, 900, 300),
			Color:  chooseControlColor(),
			Entity: Entity{UpdateFunc: func(time.Duration) {}},
		}, OnEvent: func() {}},
	}

	var testRects []*Button
	testRects = makeTestRects(controlRects, []pixel.RGBA{}, &testRects)

	saveButton := Button{
		ColoredRect: ColoredRect{
			Bounds: pixel.R(400, 400, 700, 600),
			Color:  pixel.RGB(0.8, 0.8, 0.8),
			Entity: Entity{UpdateFunc: func(time.Duration) {}},
		},
		OnEvent: func() {
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

	frameTimes := []time.Time{time.Now()}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 100), basicAtlas)

	clickIndicator := Button{ColoredRect: ColoredRect{Bounds: win.Bounds(), Color: pixel.RGB(0, 0, 0), Entity: Entity{surface: imdraw.New(nil)}}}
	clickIndicator.UpdateFunc = func(time.Duration) {
		clickIndicator.Color.G = 0
		clickIndicator.Color.A = 0
		clickIndicator.Bounds = win.Bounds()
	}
	clickIndicator.OnEvent = func() {
		clickIndicator.Color.G = 1
		clickIndicator.Color.A = 1
		clickIndicator.Bounds = pixel.R(0, 10, 10, 20)
	}

	collisionIndicator := Button{
		ColoredRect: ColoredRect{
			Bounds: pixel.R(0, 0, 10, 10),
			Color:  pixel.RGB(0, 1, 0),
			Entity: Entity{UpdateFunc: func(time.Duration) {}},
		},
		OnEvent: func() {},
	}

	entities := []*Button{}
	entities = append(entities, &controlRects[0], &controlRects[1])
	entities = append(entities, testRects...)
	entities = append(entities, &saveButton, &clickIndicator, &collisionIndicator)

	clicked := Event{}

	collisionIndicator.OnEvent = func() {
		for i, e := range entities {
			if e.Contains(clicked.mousePos) {
				collisionIndicator.Bounds = pixel.R(float64(10*i), 0, float64(10*i+10), 10)
				collisionIndicator.Color = e.Color
				collisionIndicator.Color.A = 1
			}
		}
	}
	collisionIndicator.UpdateFunc = func(time.Duration) {
		collisionIndicator.Bounds = win.Bounds()
		collisionIndicator.Color = pixel.RGB(0, 0, 0)
		collisionIndicator.Color.A = 0
	}

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		basicTxt.Clear()

		// TODO: make fps counter an entity
		frameTimes = append(frameTimes, time.Now())
		for frameTimes[len(frameTimes)-1].Sub(frameTimes[0]).Seconds() >= 1 {
			frameTimes = frameTimes[1:]
		}

		click := win.JustPressed(pixelgl.MouseButton1)
		if click {
			clicked.mousePos = win.MousePosition()
		}
		for _, e := range entities {
			e.Update(frameTimes[len(frameTimes)-1].Sub(frameTimes[len(frameTimes)-2]))
			if click {
				e.Handle(clicked)
			}
			e.Draw(win)
		}

		// Draw FPS tracker
		fmt.Fprintln(basicTxt, len(frameTimes))

		// Draw step counter
		// TODO: actually use formatting
		fmt.Fprint(basicTxt, "Step ")
		fmt.Fprintln(basicTxt, len(chosenTestColors)%7)
		basicTxt.Draw(win, pixel.IM)

		win.Update()
	}
}

func chooseControlColor() pixel.RGBA {
	return pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
}

func firstChooseTestColors(ctc []pixel.RGBA) (out []pixel.RGBA) {
	length := len(ctc)
	step := length % 7
	base := length - step

	inColor := pixel.RGB(0.5, 0.5, 0.5)
	if base > 0 {
		inColor = ctc[base-1]
	}

	if step < 4 {
		// Over the course of 4 steps, this will generate all 8 colors at a given distance from the input color.
		// The offset will start out as 1/3 and cut in half every 7 steps, i.e. every repetition of the tournament bracket.
		offset := 1 / (float64(int(3) << (base / 7)))

		// Expect R, G, and B values to occasionally exceed 1 or fall below 0.
		// TODO: clamp that.

		// Choose the first color: per step, G and R offsets respectively follow the pattern --, -+, +-, ++
		out = []pixel.RGBA{pixel.RGB(inColor.R-offset, inColor.G-offset, inColor.B-offset)}
		if step&1 == 1 {
			out[0].R += 2 * offset
		}
		if step&2 == 2 {
			out[0].G += 2 * offset
		}

		// Choose the second color: same as the first, but switch B offset from - to +.
		// After 4 steps, this generates all 8 sequences of three + or - offsets.
		out = append(out, out[0])
		out[1].B += 2 * offset
	} else {
		// Semifinals round 1: when step == 4, output ctc[base] and ctc[1+base]
		// Semifinals round 2: when step == 5, output ctc[2+base] and ctc[3+base]
		// Finals: when step == 6, output ctc[4+base] and ctc[5+base]
		out = []pixel.RGBA{ctc[step+length-8], ctc[step+length-7]}
	}

	return out
}

func makeTestRects(controlRects []Button, ctc []pixel.RGBA, testRects *[]*Button) (tr []*Button) {
	// Generate a new set of colors to compare
	for idx, col := range firstChooseTestColors(ctc) {
		// Make the rects, which when clicked call this function again
		rect := Button{ColoredRect: ColoredRect{Bounds: pixel.R(500, 100+float64(idx*100), 600, 200+float64(idx*100)), Color: col, Entity: Entity{UpdateFunc: func(time.Duration) {}}}}
		rect.OnEvent = func() {
			chosenTestColors = append(chosenTestColors, rect.Color)

			a := makeTestRects(controlRects, chosenTestColors, testRects)
			for len(*testRects) < 2 {
				*testRects = append(*testRects, &Button{})
			}
			*(*testRects)[0] = *a[0]
			*(*testRects)[1] = *a[1]
		}
		tr = append(tr, &rect)
	}
	return tr
}
