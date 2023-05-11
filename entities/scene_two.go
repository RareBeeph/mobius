package entities

import (
	"colorspacer/types"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	*Scene1 = *Scene
	Scene2.Children = []types.E{&SceneReturnButton, &ClickIndicator, &CollisionIndicator, &FpsC}

	*Scene = *Scene2

	SceneReturnButton.OnEvent = func() {
		*Scene = *Scene1
	}
}

var Scene1 = &types.Entity{}
var Scene2 = &types.Entity{}

var SceneReturnButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(800, 450, 950, 550), Color: pixel.RGB(0.6, 0.6, 0.6)},
}
