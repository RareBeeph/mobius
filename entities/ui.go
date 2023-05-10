package entities

import (
	"colorspacer/db/model"
	"colorspacer/db/query"
	"colorspacer/types"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func init() {
	AllEntities = &[]types.CR{}
}

func Initialize(win *pixelgl.Window, clicked *types.Event) {
	// Specifically not init because pixelgl needs to Run() before this

	ClickIndicator.Bounds = win.Bounds()
	ClickIndicator.UpdateFunc = func(time.Duration) {
		ClickIndicator.Color.G = 0
		ClickIndicator.Color.A = 0
		ClickIndicator.Bounds = win.Bounds()
	}
	ClickIndicator.OnEvent = func() {
		ClickIndicator.Color.G = 1
		ClickIndicator.Color.A = 1
		ClickIndicator.Bounds = pixel.R(0, 10, 10, 20)
	}

	*AllEntities = []types.CR{&SaveButton, &ClickIndicator, &CollisionIndicator, &ControlRects[0], &ControlRects[1], TestRects[0], TestRects[1], &SceneButton}

	CollisionIndicator.OnEvent = func() {
		for i, e := range *AllEntities {
			if e.Contains(clicked.MousePos) {
				CollisionIndicator.Bounds = pixel.R(float64(10*i), 0, float64(10*i+10), 10)
				CollisionIndicator.Color = e.GetColor()
				CollisionIndicator.Color.A = 1
			}
		}
	}
	CollisionIndicator.UpdateFunc = func(time.Duration) {
		CollisionIndicator.Bounds = win.Bounds()
		CollisionIndicator.Color = pixel.RGB(0, 0, 0)
		CollisionIndicator.Color.A = 0
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	FpsC.Text = text.New(pixel.V(100, 100), basicAtlas)
	FpsC.FrameTimes = []time.Time{time.Now()}

	FpsC.UpdateFunc = func(dt time.Duration) {
		FpsC.Text.Clear()

		FpsC.FrameTimes = append(FpsC.FrameTimes, FpsC.FrameTimes[len(FpsC.FrameTimes)-1].Add(dt)) // Not strictly synced to the time kept track of in main
		for FpsC.FrameTimes[len(FpsC.FrameTimes)-1].Sub(FpsC.FrameTimes[0]).Seconds() >= 1 {
			FpsC.FrameTimes = FpsC.FrameTimes[1:]
		}

		FpsC.StepCount = len(ChosenTestColors) % 7
	}

	AllTexts = []*types.FpsCounter{&FpsC}

	copy(Graph.BasisMatrix[:], types.DefaultBasisMatrix[:])
	Graph.UpdateFunc = func(dt time.Duration) {
		if len(ChosenTestColors) > 0 {
			Graph.Curve.ControlPoints = []pixel.RGBA{ControlRects[0].Color, ChosenTestColors[len(ChosenTestColors)-1], ControlRects[1].Color}
		}
	}

	SceneButton.OnEvent = func() {
		InitSceneTwo(win, clicked)
	}
}

var AllEntities *[]types.CR
var AllTexts []*types.FpsCounter
var Graph = &types.CurveDisplay{
	Center: pixel.V(150, 500),
	Curve:  types.RgbCurve{ControlPoints: []pixel.RGBA{ControlRects[0].Color, ControlRects[1].Color}},
	Bounds: pixel.R(50, 400, 250, 600),
}

var ClickIndicator = types.Button{ColoredRect: types.ColoredRect{Color: pixel.RGB(0, 0, 0)}}
var CollisionIndicator = types.Button{ColoredRect: types.ColoredRect{Bounds: pixel.R(0, 0, 10, 10), Color: pixel.RGB(0, 1, 0)}}

var FpsC types.FpsCounter

var SaveButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(400, 400, 700, 600), Color: pixel.RGB(0.8, 0.8, 0.8)},
	OnEvent: func() {
		m := query.Midpoint

		if len(ChosenTestColors) > 0 {
			start := model.NewColorFromRgba(ControlRects[0].Color)
			end := model.NewColorFromRgba(ControlRects[1].Color)
			mid := model.NewColorFromRgba(ChosenTestColors[len(ChosenTestColors)-1])
			a := model.Midpoint{
				Startpoint: *start,
				Endpoint:   *end,
				Midpoint:   *mid,
			}
			m.Create(&a)

			//Debug
			b, _ := m.Preload(m.Midpoint).Last()
			log.Printf("ID: %d, R: %f, G: %f, B: %f", b.Midpoint.ID, b.Midpoint.R, b.Midpoint.G, b.Midpoint.B)
		}
	},
}

var SceneButton = types.Button{
	ColoredRect: types.ColoredRect{Bounds: pixel.R(800, 450, 950, 550), Color: pixel.RGB(0.6, 0.6, 0.6)},
}
