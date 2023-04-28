package types

type Button struct {
	ColoredRect

	OnEvent func()
}

func (b *Button) Handle(event Event) {
	if b.Contains(event.MousePos) {
		b.OnEvent()
	}
}
