package entities

import (
	"colorspacer/types"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	SceneReturnButton.OnEvent = func() {
		Initialize(win, clicked)
	}

	*AllEntities = []types.CR{&SceneReturnButton, &ClickIndicator, &CollisionIndicator}
}

var SceneReturnButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(800, 450, 950, 550), Color: pixel.RGB(0.6, 0.6, 0.6)},
}
