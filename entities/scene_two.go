package entities

import (
	"colorspacer/db/model"
	"colorspacer/db/query"
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
	Scene2.Children = []types.E{&SceneReturnButton, &ClickIndicator, &CollisionIndicator, &S2ControlColor, &S2TestColor, &S2Slider, &S2Control2, &ProgressButton, &S2Control3, &MetricLogger, &MetricSaveButton, &GraphSlider, &MetricGraph, &FpsC}

	*Scene = *Scene2

	SceneReturnButton.OnEvent = func(e *types.Event) {
		*Scene = *Scene1
	}
}

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	metric[0][0] = 1
	lengths[0][0] = 1

	S2Slider.UpdateFunc = func(dt time.Duration) {
		S2Slider.Color = S2TestColor.Color
	}

	copy(MetricGraph.BasisMatrix[:], types.DefaultBasisMatrix[:])

	measureMetric(0, 0, 0) // Sets S2Control2.Color, S2Control3.Color, and ProgressButton.OnEvent
}

var sceneTwoInitialized = false
var Scene1 = &types.Entity{}
var Scene2 = &types.Entity{}

var SceneReturnButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(800, 450, 950, 550), Color: pixel.RGB(0.6, 0.6, 0.6)},
	Label:       "Return to scene 1",
}

var S2TestColor = types.ColoredRect{
	Bounds: pixel.R(400, 200, 500, 300),
	Color:  S2ControlColor.Color,
}

var coloroffset = 0.05 // Higher means more proportionally reliable measurements (in theory), but a worse approximation of the tangent space

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
			Color:  S2ControlColor.Color,
			Bounds: pixel.R(80+S2ControlColor.Color.R*600, 130, 120+S2Control2.Color.R*600, 170),
		},
		EventTypeHandled: types.Drag,
	},
	ClampMin:    100,
	ClampMax:    700,
	TargetValue: &S2TestColor.Color.R,
	OutputMin:   0,
	OutputMax:   1,
}

var ProgressButton = types.Button{
	ColoredRect: types.ColoredRect{
		Bounds: pixel.R(350, 500, 550, 600),
		Color:  pixel.RGB(0.8, 0.8, 0.8),
	},
	Label: "Next measurement step",
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
	Label: "Print metric to log",
}

var MetricGraph = types.MetricDisplay{
	CurveDisplay: types.CurveDisplay{
		Center: pixel.V(150, 500),
		Bounds: pixel.R(50, 400, 250, 600),
	},
	ColorOffset:     coloroffset,
	CenterCol:       S2ControlColor.Color,
	CenterDepth:     375,
	ThicknessFactor: 2.5,
}

var GraphSlider = types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color:  pixel.RGB(0.5, 0.5, 0.5),
			Bounds: pixel.R(130, 680, 170, 720),
		},
	},
	ClampMin:    50,
	ClampMax:    250,
	TargetValue: &MetricGraph.CenterDepth,
	OutputMin:   500,
	OutputMax:   150,
}

var MetricSaveButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(600, 200, 900, 400), Color: pixel.RGB(0.8, 0.8, 0.8)},
	OnEvent: func(e *types.Event) {
		m := query.Metric

		a := model.NewMetricFromArray(metric)
		m.Create(&a)

		//Debug
		b, _ := m.Last()
		log.Printf("ID: %d, RR: %f, GG: %f, BB: %f, RG: %f, RB: %f, GB: %f", b.ID, b.RedSquared, b.GreenSquared, b.BlueSquared, b.RedDotGreen, b.RedDotBlue, b.GreenDotBlue)
	},
	Label: "Save metric to database",
}

var metric [3][3]float64
