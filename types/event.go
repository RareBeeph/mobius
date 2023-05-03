package types

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Event struct {
	MousePos   pixel.Vec
	InitialPos pixel.Vec
	Buttons    []pixelgl.Button
}

func (e *Event) Contains(b pixelgl.Button) bool {
	for _, t := range e.Buttons {
		if b == t {
			return true
		}
	}
	return false
}
