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

type E interface {
	Update(time.Duration)
	Draw(*pixelgl.Window)
	Handle(Event)
	Handles() bool
}

func (entity *Entity) Update(deltatime time.Duration) {
	if entity.UpdateFunc != nil {
		entity.UpdateFunc(deltatime)
	}
}

func (entity *Entity) Draw(window *pixelgl.Window) {
	entity.surface.Draw(window)
}

func (entity *Entity) Handle(event Event) {

}

func (entity *Entity) Handles() bool {
	return false
}
