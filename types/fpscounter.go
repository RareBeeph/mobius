package types

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type FpsCounter struct {
	Entity

	StepCount  int
	FrameTimes []time.Time
	Text       *text.Text
}

func NewFpsCounter(pos pixel.Vec) (fpsc *FpsCounter) {
	fpsc = &FpsCounter{}
	fpsc.FrameTimes = []time.Time{time.Now()}
	fpsc.Text = text.New(pos, text.NewAtlas(basicfont.Face7x13, text.ASCII))

	return fpsc
}

func (fpsc *FpsCounter) Update(dt time.Duration) {
	fpsc.FrameTimes = append(fpsc.FrameTimes, fpsc.FrameTimes[len(fpsc.FrameTimes)-1].Add(dt)) // Not strictly synced to the time kept track of in main
	for fpsc.FrameTimes[len(fpsc.FrameTimes)-1].Sub(fpsc.FrameTimes[0]).Seconds() >= 1 {
		fpsc.FrameTimes = fpsc.FrameTimes[1:]
	}

	if fpsc.UpdateFunc != nil {
		fpsc.UpdateFunc(dt)
	}
}

func (fpsc *FpsCounter) Draw(window *pixelgl.Window) {
	fpsc.Text.Clear()

	fmt.Fprintln(fpsc.Text, len(fpsc.FrameTimes))

	fmt.Fprint(fpsc.Text, "Step ")
	fmt.Fprintln(fpsc.Text, fpsc.StepCount)

	fpsc.Text.Draw(window, pixel.IM)
}
