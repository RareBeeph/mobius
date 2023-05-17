package entities

import (
	"colorspacer/types"
	"log"
	"math"
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
	Scene2.Children = []types.E{&SceneReturnButton, &ClickIndicator, &CollisionIndicator, &S2ControlColor, &S2TestColor, &S2Slider, &S2Control2, &ProgressButton, &S2Control3, &MetricLogger, &MetricGraph, &FpsC}

	*Scene = *Scene2

	SceneReturnButton.OnEvent = func(e *types.Event) {
		*Scene = *Scene1
	}
}

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	metric[0][0] = 1
	lengths[0][0] = 1

	sliderRange := S2Slider.ClampMax - S2Slider.ClampMin
	centerPos := S2Slider.ClampMin + S2Slider.Color.R*sliderRange

	S2Slider.Bounds = pixel.R(centerPos-20, 130, centerPos+20, 170)

	S2Slider.OnEvent = func(e *types.Event) {
		S2Slider.Bounds.Min.X = e.MousePos.X - 20
		S2Slider.Bounds.Max.X = e.MousePos.X + 20
		S2Slider.Clamp()

		colorAmount := (S2Slider.Bounds.Center().X - S2Slider.ClampMin) / sliderRange
		S2Slider.Color.R = colorAmount
	}

	S2TestColor.UpdateFunc = func(dt time.Duration) {
		S2TestColor.Color = S2Slider.Color
	}

	MetricGraph.CenterCol = S2ControlColor.Color
	copy(MetricGraph.BasisMatrix[:], types.DefaultBasisMatrix[:])

	measureMetric(1, 1, 0)
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
}

var coloroffset = 0.1 // Higher means more proportionally reliable measurements (in theory), but a worse approximation of the tangent space

var S2ControlColor = types.ColoredRect{
	Bounds: pixel.R(300, 200, 400, 300),
	Color:  chooseControlColor().Scaled(0.8).Add(pixel.RGB(0.1, 0.1, 0.1)),
}

var S2Control2 = types.ColoredRect{
	Bounds: pixel.R(400, 300, 500, 400),
	Color:  S2Control3.Color.Add(pixel.RGB(0, coloroffset, 0)), // Note: unclamped, and should actually be a constant
}

var S2Control3 = types.ColoredRect{
	Bounds: pixel.R(300, 300, 400, 400),
	Color:  S2ControlColor.Color, // Should actually be a constant
}

var S2Slider = types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color: S2ControlColor.Color,
		},
		EventTypeHandled: types.Drag,
	},
	ClampMin: 100,
	ClampMax: 700,
}

var ProgressButton = types.Button{
	ColoredRect: types.ColoredRect{
		Bounds: pixel.R(350, 500, 550, 600),
		Color:  pixel.RGB(0.8, 0.8, 0.8),
	},
}

var MetricLogger = types.Button{
	ColoredRect: types.ColoredRect{
		Bounds: pixel.R(600, 500, 760, 600),
		Color:  pixel.RGB(0.8, 0.8, 0.8),
	},
	OnEvent: func(e *types.Event) {
		log.Print("Metric: ")
		log.Println(metric)

		var modifiedAngles [3]float64
		for i := range angles {
			modifiedAngles[i] = angles[i] * 180 / math.Pi
		}
		log.Print("Angles (degrees): ")
		log.Println(modifiedAngles)
	},
}

var MetricGraph = types.MetricDisplay{
	CurveDisplay: types.CurveDisplay{
		Center: pixel.V(150, 500),
		Bounds: pixel.R(50, 400, 250, 600),
	},
	ColorOffset: coloroffset,
}

var metric [3][3]float64
