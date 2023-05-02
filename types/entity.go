package types

import (
	"time"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Entity struct {
	surface *imdraw.IMDraw

	UpdateFunc func(time.Duration)
}

type EventHandler interface {
	Update(time.Duration)
	Draw(*pixelgl.Window)
	Handle(Event)
	Handles() bool
}

type E = EventHandler

func (entity *Entity) GuardSurface() {
	// Generate new surface if we were not provided one
	if entity.surface == nil {
		entity.surface = imdraw.New(nil)
	}

	entity.surface.Clear()
}

func (entity *Entity) Update(deltatime time.Duration) {
	if entity.UpdateFunc != nil {
		entity.UpdateFunc(deltatime)
	}
}

func (entity *Entity) Draw(window *pixelgl.Window) {
	entity.surface.Draw(window) // As of writing this comment, this is never run. It should crash (null reference) if it were.
}

func (entity *Entity) Handle(event Event) {

}

func (entity *Entity) Handles() bool {
	return false
}
