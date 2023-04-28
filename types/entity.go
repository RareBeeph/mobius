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

/*
Interfaces are usually named so that they fit in a phrase like
"object implements the EventHandler interface" and so I would
recommend renaming it. HOWEVER...
*/
type EventHandler interface {
	Update(time.Duration)
	Draw(*pixelgl.Window)
	Handle(Event)
	Handles() bool
}

/*
...the code ergonomics of having it named E are undeniable.
We get the best of both worlds by using type aliases, which
are declared like so:
*/
type E = EventHandler

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
