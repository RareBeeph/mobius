package types

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Slider struct {
	Button

	InitialBounds pixel.Rect
	ClampMin      float64
	ClampMax      float64
}

func (s *Slider) Contained(point pixel.Vec) bool {
	return (point.X >= s.InitialBounds.Min.X &&
		point.X < s.InitialBounds.Max.X &&
		point.Y >= s.InitialBounds.Min.Y &&
		point.Y < s.InitialBounds.Max.Y)
}

func (s *Slider) Handle(event *Event) {
	s.OnEvent(event)
}

func (s *Slider) Handles(event *Event) bool {
	if event.EventType == Click && event.Contains(pixelgl.MouseButton1) {
		s.InitialBounds = s.Bounds
		return false
	}
	if s.Contained(event.InitialPos) && event.Contains(pixelgl.MouseButton1) {
		return true
	}
	return false
}

func (s *Slider) Clamp() {
	if s.Bounds.Center().X > s.ClampMax {
		dx := s.Bounds.Center().X - s.ClampMax
		s.Bounds.Max.X -= dx
		s.Bounds.Min.X -= dx
	}
	if s.Bounds.Center().X < s.ClampMin {
		dx := s.Bounds.Center().X - s.ClampMin
		s.Bounds.Max.X -= dx
		s.Bounds.Min.X -= dx
	}
}
