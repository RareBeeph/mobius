package main

import (
	"colorspacer/types"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type Dispatch struct {
	Buttons []types.CR // Collision indicator needed to find the color of all buttons, so this can't be []types.E
}

// These feel a bit repetitive
func (dispatch *Dispatch) Update(deltatime time.Duration) {
	for _, e := range dispatch.Buttons {
		e.Update(deltatime)
	}
}

func (dispatch *Dispatch) Handle(event types.Event) {
	for _, e := range dispatch.Buttons {
		e.Handle(event)
	}
}

func (dispatch *Dispatch) Draw(win *pixelgl.Window) {
	for _, e := range dispatch.Buttons {
		e.Draw(win)
	}
}
