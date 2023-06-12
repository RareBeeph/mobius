package entities

import (
	"colorspacer/types"
	"time"

	"github.com/faiface/pixel"
)

/* TODO:
choose 2 grays
sliderify a gray
repeat
interpret the results to build a standard measure of length

choose 2 colors
sliderify a gray
interpret the results in terms of the standard
repeat

interpret the results to build a model of distances between un-measured colors as well
*/

var SceneGrays = &types.Entity{}
var SGinitialized = false

func InitSceneGrays() {
	SGTestColor.(*types.ColoredRect).UpdateFunc = func(dt time.Duration) {
		SGTestColor.(*types.ColoredRect).Color = pixel.RGB(testGray, testGray, testGray)
	}

	SGSlider.(*types.Slider).UpdateFunc = func(dt time.Duration) {
		SGSlider.(*types.Slider).Color = pixel.RGB(testGray, testGray, testGray)
	}
}

var SGTestColor = SceneGrays.AddChild(&types.ColoredRect{
	Bounds: pixel.R(400, 200, 500, 300),
	Color:  ControlGray,
})

var SGControlColor = SceneGrays.AddChild(&types.ColoredRect{
	Bounds: pixel.R(300, 200, 400, 300),
	Color:  pixel.RGB(0, 0, 0),
})

var SGControl2 = SceneGrays.AddChild(&types.ColoredRect{
	Bounds: pixel.R(400, 300, 500, 400),
	Color:  ControlGray.Add(pixel.RGB(grayoffset, grayoffset, grayoffset)), // Note: unclamped, and should actually be a constant
})

var SGControl3 = SceneGrays.AddChild(&types.ColoredRect{
	Bounds: pixel.R(300, 300, 400, 400),
	Color:  ControlGray,
})

var ControlGray pixel.RGBA
var grayoffset = 0.10
var testGray = 0.5

var SGSlider = SceneGrays.AddChild(&types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color:  pixel.RGB(testGray, testGray, testGray),
			Bounds: pixel.R(80+testGray*600, 130, 120+testGray*600, 170),
		},
		EventTypeHandled: types.Drag,
	},
	ClampMin: 100,
	ClampMax: 700,
	// Should be multiple target values; currently unsupported
	TargetValue: &testGray,
	OutputMin:   0,
	OutputMax:   1,
})
