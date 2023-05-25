package types

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ColoredRect struct {
	Entity

	Bounds pixel.Rect
	Color  pixel.RGBA
}

type ColoredRectHandler interface {
	EI

	Contains(pixel.Vec) bool
	GetColor() pixel.RGBA
}

type CRI = ColoredRectHandler

func (r *ColoredRect) Contains(point pixel.Vec) bool {
	return (point.X >= r.Bounds.Min.X &&
		point.X < r.Bounds.Max.X &&
		point.Y >= r.Bounds.Min.Y &&
		point.Y < r.Bounds.Max.Y)
}

func (r *ColoredRect) Draw(window *pixelgl.Window) {
	r.GuardSurface()

	r.surface.Color = r.Color
	r.surface.Push(r.Bounds.Min)
	r.surface.Push(r.Bounds.Max)
	r.surface.Rectangle(0)

	r.surface.Draw(window)
}

func (r *ColoredRect) GetColor() pixel.RGBA {
	return r.Color
}
