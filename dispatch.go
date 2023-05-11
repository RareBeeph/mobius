package main

import (
	"colorspacer/types"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type Dispatch struct {
	Entities *types.Entities
}

func (dispatch *Dispatch) Update(deltatime time.Duration) {
	for _, e := range *dispatch.Entities {
		e.Update(deltatime)
	}
}

func (dispatch *Dispatch) Handle(event *types.Event) {
	for _, e := range *dispatch.Entities {
		types.Receive(e, event)
	}
}

func (dispatch *Dispatch) Draw(win *pixelgl.Window) {
	for _, e := range *dispatch.Entities {
		e.Draw(win)
	}
}
