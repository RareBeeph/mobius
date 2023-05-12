package entities

import (
	"colorspacer/types"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func SwitchToSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	if !sceneTwoInitialized {
		InitSceneTwo(win, clicked)
		sceneTwoInitialized = true
	}

	*Scene1 = *Scene
	Scene2.Children = []types.E{&SceneReturnButton, &ClickIndicator, &CollisionIndicator, &S2ControlColor, &S2TestColor, &S2Slider, &S2Control2, &FpsC}

	*Scene = *Scene2

	SceneReturnButton.OnEvent = func(e *types.Event) {
		*Scene = *Scene1
	}
}

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	sliderRange := S2Slider.ClampMax - S2Slider.ClampMin
	centerPos := S2Slider.ClampMin + S2Slider.Color.G*sliderRange

	S2Slider.Bounds = pixel.R(centerPos-20, 130, centerPos+20, 170)

	S2Slider.OnEvent = func(e *types.Event) {
		S2Slider.Bounds.Min.X = e.MousePos.X - 20
		S2Slider.Bounds.Max.X = e.MousePos.X + 20
		S2Slider.Clamp()

		S2Slider.Color.G = (S2Slider.Bounds.Center().X - S2Slider.ClampMin) / sliderRange
	}

	S2TestColor.UpdateFunc = func(dt time.Duration) {
		S2TestColor.Color = S2Slider.Color
	}
}

var sceneTwoInitialized = false
var Scene1 = &types.Entity{}
var Scene2 = &types.Entity{}

var SceneReturnButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(800, 450, 950, 550), Color: pixel.RGB(0.6, 0.6, 0.6)},
}

var S2TestColor = types.ColoredRect{
	Bounds: pixel.R(400, 200, 500, 300),
	Color:  S2ControlColor.Color,
	Entity: types.Entity{
		UpdateFunc: func(dt time.Duration) {

		}},
}

var S2ControlColor = types.ColoredRect{
	Bounds: pixel.R(300, 200, 400, 400),
	Color:  chooseControlColor(),
}

var S2Control2 = types.ColoredRect{
	Bounds: pixel.R(400, 300, 500, 400),
	Color:  S2ControlColor.Color.Add(pixel.RGB(0.05, 0, 0)), // Note: unclamped
}

var S2Slider = types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color: S2ControlColor.Color,
		},
		EventTypeHandled: types.Drag,
	},
	ClampMin: 300,
	ClampMax: 500,
}
