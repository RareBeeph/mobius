package main

type Button struct {
	ColoredRect

	procedure func()
}

func (b *Button) Handle(event Event) {
	if b.Contains(event.mousePos) {
		b.procedure()
	}
}
