package entities

import (
	"colorspacer/types"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

func init() {
	TestRects = makeTestRects(ControlRects, []pixel.RGBA{}, &TestRects)
}

var ControlRects = []types.ColoredRect{
	{Bounds: pixel.R(200, 100, 500, 300), Color: chooseControlColor()},
	{Bounds: pixel.R(600, 100, 900, 300), Color: chooseControlColor()},
}

var ChosenTestColors []pixel.RGBA

var TestRects []*types.Button

func chooseControlColor() pixel.RGBA {
	rand.Seed(time.Now().UnixMicro()) // Temp until I figure out why rand isn't being random without this here
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

func makeTestRects(controlRects []types.ColoredRect, ctc []pixel.RGBA, testRects *[]*types.Button) (tr []*types.Button) {
	// Generate a new set of colors to compare
	for idx, col := range firstChooseTestColors(ctc) {
		// Make the rects, which when clicked call this function again
		rect := types.Button{ColoredRect: types.ColoredRect{Bounds: pixel.R(500, 100+float64(idx*100), 600, 200+float64(idx*100)), Color: col}}
		rect.OnEvent = func() {
			ChosenTestColors = append(ChosenTestColors, rect.Color)

			a := makeTestRects(controlRects, ChosenTestColors, testRects)
			for len(*testRects) < 2 {
				*testRects = append(*testRects, &types.Button{})
			}
			*(*testRects)[0] = *a[0]
			*(*testRects)[1] = *a[1]
		}
		tr = append(tr, &rect)
	}
	return tr
}
