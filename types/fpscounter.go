package types

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type FpsCounter struct {
	Entity

	StepCount  int
	FrameTimes []time.Time
	Text       *text.Text
}

func (fpsc *FpsCounter) Draw(window *pixelgl.Window) {
	fmt.Fprintln(fpsc.Text, len(fpsc.FrameTimes))

	fmt.Fprint(fpsc.Text, "Step ")
	fmt.Fprintln(fpsc.Text, fpsc.StepCount)
	fpsc.Text.Draw(window, pixel.IM)

	fpsc.Text.Draw(window, pixel.IM)
}
