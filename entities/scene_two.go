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
	*Scene = *Scene2
}

func InitSceneTwo(win *pixelgl.Window, clicked *types.Event) {
	metric[0][0] = 1
	lengths[0][0] = 1

	S2Slider.(*types.Slider).UpdateFunc = func(dt time.Duration) {
		S2Slider.(*types.Slider).Color = S2TestColor.(*types.ColoredRect).Color
	}

	allButtons := Scene2.FindAllChildren(func(ei types.EI) bool {
		_, ok := ei.(types.CRI)
		return ok
	})
	// Mostly copied from ui.go
	// TODO: make indicators their own type
	// Occasionally drops inputs for unknown reason
	s2col := S2CollisionIndicator.(*types.Button)
	s2col.OnEvent = func(e *types.Event) {
		for i, e := range allButtons {
			if e.(types.CRI).Contains(clicked.MousePos) {
				s2col.Bounds = pixel.R(float64(10*i), 0, float64(10*i+10), 10)
				s2col.Color = e.(types.CRI).GetColor()
				s2col.Color.A = 1
			}
		}
	}
	s2col.UpdateFunc = func(time.Duration) {
		s2col.Bounds = win.Bounds()
		s2col.Color = pixel.RGB(0, 0, 0)
		s2col.Color.A = 0
	}

	// Mostly copied from ui.go
	s2clk := S2ClickIndicator.(*types.Button)
	s2clk.Bounds = win.Bounds()
	s2clk.UpdateFunc = func(time.Duration) {
		s2clk.Color.G = 0
		s2clk.Color.A = 0
		s2clk.Bounds = win.Bounds()
	}
	s2clk.OnEvent = func(e *types.Event) {
		s2clk.Color.G = 1
		s2clk.Color.A = 1
		s2clk.Bounds = pixel.R(0, 10, 10, 20)
	}

	copy(MetricGraph.(*types.MetricDisplay).BasisMatrix[:], types.DefaultBasisMatrix[:])

	measureMetric(0, 0, 0) // Sets S2Control2.Color, S2Control3.Color, and ProgressButton.OnEvent
}

var sceneTwoInitialized = false
var Scene1 = &types.Entity{}
var Scene2 = &types.Entity{}

var ControlColor = chooseControlColor().Scaled(0.8).Add(pixel.RGB(0.1, 0.1, 0.1))

var SceneReturnButton = Scene2.AddChild(&types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(800, 450, 950, 550), Color: pixel.RGB(0.6, 0.6, 0.6)},
	Label:       "Return to scene 1",
	OnEvent: func(e *types.Event) {
		*Scene = *Scene1
	},
})

var S2TestColor = Scene2.AddChild(&types.ColoredRect{
	Bounds: pixel.R(400, 200, 500, 300),
	Color:  ControlColor,
})

var coloroffset = 0.05 // Higher means more proportionally reliable measurements (in theory), but a worse approximation of the tangent space

var S2ControlColor = Scene2.AddChild(&types.ColoredRect{
	Bounds: pixel.R(300, 200, 400, 300),
	Color:  ControlColor,
})

var S2Control2 = Scene2.AddChild(&types.ColoredRect{
	Bounds: pixel.R(400, 300, 500, 400),
	Color:  ControlColor.Add(pixel.RGB(0, coloroffset, 0)), // Note: unclamped, and should actually be a constant
})

var S2Control3 = Scene2.AddChild(&types.ColoredRect{
	Bounds: pixel.R(300, 300, 400, 400),
	Color:  ControlColor, // Should actually be a constant
})

var S2Slider = Scene2.AddChild(&types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color:  ControlColor,
			Bounds: pixel.R(80+ControlColor.R*600, 130, 120+ControlColor.R*600, 170),
		},
		EventTypeHandled: types.Drag,
	},
	ClampMin:    100,
	ClampMax:    700,
	TargetValue: &S2TestColor.(*types.ColoredRect).Color.R,
	OutputMin:   0,
	OutputMax:   1,
})

var ProgressButton = Scene2.AddChild(&types.Button{
	ColoredRect: types.ColoredRect{
		Bounds: pixel.R(350, 500, 550, 600),
		Color:  pixel.RGB(0.8, 0.8, 0.8),
	},
	Label: "Next measurement step",
})

var MetricLogger = Scene2.AddChild(&types.Button{
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
})

var MetricGraph = Scene2.AddChild(&types.MetricDisplay{
	CurveDisplay: types.CurveDisplay{
		Center: pixel.V(150, 500),
		Bounds: pixel.R(50, 400, 250, 600),
	},
	ColorOffset:     coloroffset,
	CenterCol:       ControlColor,
	CenterDepth:     375,
	ThicknessFactor: 2.5,
})

var GraphSlider = Scene2.AddChild(&types.Slider{
	Button: types.Button{
		ColoredRect: types.ColoredRect{
			Color:  pixel.RGB(0.5, 0.5, 0.5),
			Bounds: pixel.R(130, 680, 170, 720),
		},
	},
	ClampMin:    50,
	ClampMax:    250,
	TargetValue: &MetricGraph.(*types.MetricDisplay).CenterDepth,
	OutputMin:   500,
	OutputMax:   150,
})

var MetricSaveButton = Scene2.AddChild(&types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(600, 200, 900, 400), Color: pixel.RGB(0.8, 0.8, 0.8)},
	OnEvent: func(e *types.Event) {
		m := query.Metric

		a := model.NewMetricFromArray(metric, *model.NewColorFromRgba(ControlColor))

		// Slightly janky (float64 equality checking!!!) hack to make entries with non-unique control colors update instead of create
		// There's probably something I could do with unique indexes or something to get the same effect, but this was easier
		info, _ := m.Where(m.ControlR.Eq(a.ControlR), m.ControlG.Eq(a.ControlG), m.ControlB.Eq(a.ControlB)).Updates(a)
		if info.RowsAffected == 0 {
			m.Save(&a)
		}

		//Debug
		b, _ := m.Last()
		log.Printf("ID: %d, RR: %f, GG: %f, BB: %f, RG: %f, RB: %f, GB: %f", b.ID, b.RedSquared, b.GreenSquared, b.BlueSquared, b.RedDotGreen, b.RedDotBlue, b.GreenDotBlue)
		log.Printf("Color:: R: %f, G: %f, B: %f", b.ControlR, b.ControlG, b.ControlB)
	},
	Label: "Save metric to database",
})

var S2CollisionIndicator = Scene2.AddChild(&types.Button{ColoredRect: types.ColoredRect{Bounds: pixel.R(0, 0, 10, 10), Color: pixel.RGB(0, 1, 0)}})

var S2ClickIndicator = Scene2.AddChild(&types.Button{ColoredRect: types.ColoredRect{Color: pixel.RGB(0, 0, 0)}})

var S2FpsCounter = Scene2.AddChild(types.NewFpsCounter(pixel.V(100, 100)))

var metric [3][3]float64
