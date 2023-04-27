package main

import (
	"time"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Entity struct {
	surface *imdraw.IMDraw

	UpdateFunc func(time.Duration)
}

func (entity *Entity) Update(deltatime time.Duration) {
	entity.UpdateFunc(deltatime)
}

func (entity *Entity) Draw(window *pixelgl.Window) {
	entity.surface.Draw(window)
}

func (entity *Entity) Handle(event Event) {

}

func Handles() bool {
	return false
}
