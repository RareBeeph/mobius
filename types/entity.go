package types

import (
	"time"

	"sync"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type EventHandlers []EI

type Entity struct {
	surface *imdraw.IMDraw
	wg      *sync.WaitGroup

	UpdateFunc func(time.Duration)
	Children   EventHandlers
}

type EventHandler interface {
	Update(time.Duration)
	Draw(*pixelgl.Window)
	Handle(*Event)
	Handles(*Event) bool
	GetWaitGroup() *sync.WaitGroup
	GetChildren() EventHandlers
	SetChildren(EventHandlers)
	GuardSurface()
}

type EI = EventHandler

func NewEntityWithParent[GenericEntity EI](parent EI, inputentity GenericEntity) GenericEntity {
	parent.SetChildren(append(parent.GetChildren(), inputentity))
	return inputentity
}

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

func Update(e EI, deltatime time.Duration) {
	e.Update(deltatime)
	for _, c := range e.GetChildren() {
		Update(c, deltatime)
	}
}

func (entity *Entity) Draw(window *pixelgl.Window) {
	entity.GuardSurface()
	entity.surface.Draw(window)
}

func Draw(e EI, window *pixelgl.Window) {
	e.GuardSurface()
	e.Draw(window)
	for _, c := range e.GetChildren() {
		Draw(c, window)
	}
}

func (entity *Entity) Handle(event *Event) {

}

func (entity *Entity) Handles(event *Event) bool {
	return false
}

func Receive(e EI, event *Event) {
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
		go func(child EI) {
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

func (e *Entity) GetChildren() EventHandlers {
	return e.Children
}

func (e *Entity) SetChildren(children EventHandlers) {
	e.Children = children
}

func AppendChild[C EventHandler](e *Entity, child C) C {
	e.SetChildren(append(e.GetChildren(), child))
	return child
}
