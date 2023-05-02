package types

import "github.com/faiface/pixel/pixelgl"

type Button struct {
	ColoredRect

	OnEvent func()
}

func (b *Button) Handle(event Event) {
	if b.Contains(event.MousePos) {
		if event.Contains(pixelgl.MouseButton1) {
			b.OnEvent()
		}
	}
}
