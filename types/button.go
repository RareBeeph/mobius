package types

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	ColoredRect

	OnEvent          func(*Event)
	EventTypeHandled EType // Defaults to 0, 'Click'
	Label            string
	Text             *text.Text
}

func (b *Button) Handle(event *Event) {
	b.OnEvent(event)
}

func (b *Button) Draw(window *pixelgl.Window) {
	b.GuardText()

	b.ColoredRect.Draw(window)

	b.Text.Clear()

	// Find the RGB-furthest color from the button's color, and write the text in that
	b.Text.Color = pixel.RGB(1-math.Round(b.Color.R), 1-math.Round(b.Color.G), 1-math.Round(b.Color.B))

	fmt.Fprint(b.Text, b.Label)
	b.Text.Draw(window, pixel.IM)
}

func (b *Button) GuardText() {
	if b.Text == nil {
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		b.Text = text.New(pixel.V(b.Bounds.Min.X+10, b.Bounds.Min.Y+10), basicAtlas)
	}
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
