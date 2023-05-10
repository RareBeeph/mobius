package types

import (
	"time"

	"sync"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Entities []Entity

type Entity struct {
	surface *imdraw.IMDraw
	wg      *sync.WaitGroup

	UpdateFunc func(time.Duration)
	Children   Entities
}

type EventHandler interface {
	Update(time.Duration)
	Draw(*pixelgl.Window)
	Handle(*Event)
	Handles(*Event) bool
	Receive(*Event)
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
	entity.GuardSurface()
	entity.surface.Draw(window) // As of writing this comment, this is never run. It should crash (null reference) if it were.
}

func (entity *Entity) Handle(event *Event) {

}

func (entity *Entity) Handles(event *Event) bool {
	return false
}

func (e *Entity) Receive(event *Event) {
	// TODO: This probably needs to be unwound so we can use the defer keyword to release
	if e.Handles(event) {
		event.Lock()
		e.Handle(event)
		event.Unlock()
	}

	// Don't echo down the tree if it's no longer needed
	if event.StopPropagating {
		return
	}

	// Otherwise, fire away
	e.wg = &sync.WaitGroup{}

	for _, c := range e.Children {
		e.wg.Add(1)
		go func(child *Entity) {
			defer e.wg.Done()
			child.Receive(event)
		}(&c)
	}

	e.wg.Wait()
}
