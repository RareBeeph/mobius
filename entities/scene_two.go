package entities

import (
	"colorspacer/types"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	copy(MetricGraph.(*types.MobiusDisplay).BasisMatrix[:], types.DefaultBasisMatrix[:])
}

var Scene2 = &types.Entity{}

var MetricGraph = Scene2.AddChild(&types.MobiusDisplay{
	Center:          pixel.V(150, 200),
	Bounds:          pixel.R(50, 100, 450, 300),
	CenterDepth:     375,
	ThicknessFactor: 4,
})

var ThetaSlider = Scene2.AddChild(&types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color:  pixel.RGB(0.5, 0.5, 0.5),
			Bounds: pixel.R(230, 350, 270, 390),
		},
	},
	TargetValue:   types.Tgrain,
	OutputMin:     5,
	OutputMax:     57,
	InitialBounds: pixel.R(230, 350, 270, 390),
	ClampMin:      50,
	ClampMax:      450,
})

var ISlider = Scene2.AddChild(&types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color:  pixel.RGB(0.5, 0.5, 0.5),
			Bounds: pixel.R(230, 300, 270, 340),
		},
	},
	TargetValue:   types.Igrain,
	OutputMin:     2,
	OutputMax:     18,
	InitialBounds: pixel.R(230, 300, 270, 340),
	ClampMin:      50,
	ClampMax:      450,
})

var S2FpsCounter = Scene2.AddChild(types.NewFpsCounter(pixel.V(100, 100)))

var metric [3][3]float64
