package types

import "github.com/faiface/pixel/pixelgl"

type Button struct {
	ColoredRect

	OnEvent func()
}

func (b *Button) Handle(event *Event) {
	if b.Handles(event) {
		b.OnEvent()
	}
}

func (b *Button) Handles(event *Event) bool {
	if b.Contains(event.MousePos) && event.Contains(pixelgl.MouseButton1) {
		return true
	}
	return false
}
