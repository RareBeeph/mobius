package types

import "github.com/faiface/pixel/pixelgl"

type Button struct {
	ColoredRect

	OnEvent          func(*Event)
	EventTypeHandled EType // Defaults to 0, 'Click'
}

func (b *Button) Handle(event *Event) {
	b.OnEvent(event)
}

func (b *Button) Handles(event *Event) bool {
	if event.EventType != b.EventTypeHandled {
		return false
	}
	if b.Contains(event.MousePos) && event.Contains(pixelgl.MouseButton1) {
		return true
	}
	return false
}
