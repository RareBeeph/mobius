package types

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type ColoredRect struct {
	Entity

	Bounds  pixel.Rect
	Color   pixel.RGBA
	Surface *imdraw.IMDraw
}

type CR interface {
	E

	Contains(pixel.Vec) bool
	GetColor() pixel.RGBA
}

func (r *ColoredRect) Contains(point pixel.Vec) bool {
	return (point.X >= r.Bounds.Min.X &&
		point.X < r.Bounds.Max.X &&
		point.Y >= r.Bounds.Min.Y &&
		point.Y < r.Bounds.Max.Y)
}

func (r *ColoredRect) Draw(window *pixelgl.Window) {
	// Generate new surface if we were not provided one
	if r.Surface == nil {
		r.Surface = imdraw.New(nil)
	}

	r.Surface.Clear()

	r.Surface.Color = r.Color
	r.Surface.Push(r.Bounds.Min)
	r.Surface.Push(r.Bounds.Max)
	r.Surface.Rectangle(0)

	r.Surface.Draw(window)
}

func (r *ColoredRect) GetColor() pixel.RGBA {
	return r.Color
}
