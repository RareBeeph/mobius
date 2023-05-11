package types

import (
	"time"

	"sync"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Entities []E

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
	GetWaitGroup() *sync.WaitGroup
	GetChildren() Entities
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

func Receive(e E, event *Event) {
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
	wg := e.GetWaitGroup()
	children := e.GetChildren()

	for _, c := range children {
		wg.Add(1)
		go func(child E) {
			defer wg.Done()
			Receive(child, event)
		}(c)
	}

	wg.Wait()
}

func (e *Entity) GetWaitGroup() *sync.WaitGroup {
	e.wg = &sync.WaitGroup{} // Temp
	return e.wg
}

func (e *Entity) GetChildren() Entities {
	return e.Children
}
