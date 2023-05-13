package entities

import (
	"colorspacer/types"
	"log"
	"math"

	"github.com/faiface/pixel"
)

func g33() {
	// Control2 is +blue
	S2Control2.Color = S2Control3.Color.Add(pixel.RGB(0, 0, coloroffset))

	ProgressButton.OnEvent = func(e *types.Event) {
		log.Println((S2Slider.Color.R - S2ControlColor.Color.R) / coloroffset)
		// One step in blue (as exemplified by control colors) == some number of steps in red; square that
		metric[2][2] = math.Pow((S2Slider.Color.R-S2ControlColor.Color.R)/coloroffset, 2) // Might not work when comparing locations
		g12()
	}
}

func g12() {
	// Control2 is +red -green
	S2Control2.Color = S2Control3.Color.Add(pixel.RGB(coloroffset, -coloroffset, 0))

	ProgressButton.OnEvent = func(e *types.Event) {
		// R dot G = -1/2 * ((R-G)^2 - R^2 - G^2)
		// R-G = (S2Slider.Color.B-S2ControlColor.Color.B)/coloroffset*sqrt(metric[2][2])
		// R^2 = metric[0][0]
		// G^2 = metric[1][1]
		log.Println((S2Slider.Color.R - S2ControlColor.Color.R) / coloroffset * math.Sqrt(metric[0][0]))
		metric[0][1] = -0.5 * (math.Pow((S2Slider.Color.R-S2ControlColor.Color.R)/coloroffset, 2)*metric[0][0] - metric[0][0] - metric[1][1])
		metric[1][0] = metric[0][1]
		g13()
	}
}

func g13() {
	// Control2 is +red -blue
	S2Control2.Color = S2Control3.Color.Add(pixel.RGB(coloroffset, 0, -coloroffset))

	ProgressButton.OnEvent = func(e *types.Event) {
		log.Println((S2Slider.Color.R - S2ControlColor.Color.R) / coloroffset * math.Sqrt(metric[0][0]))
		metric[0][2] = -0.5 * (math.Pow((S2Slider.Color.R-S2ControlColor.Color.R)/coloroffset, 2)*metric[0][0] - metric[0][0] - metric[2][2])
		metric[2][0] = metric[0][2]
		g23()
	}
}

func g23() {
	// Control2 is +green -blue
	S2Control2.Color = S2Control3.Color.Add(pixel.RGB(0, coloroffset, -coloroffset))

	ProgressButton.OnEvent = func(e *types.Event) {
		log.Println((S2Slider.Color.R - S2ControlColor.Color.R) / coloroffset * math.Sqrt(metric[0][0]))
		metric[1][2] = -0.5 * (math.Pow((S2Slider.Color.R-S2ControlColor.Color.R)/coloroffset, 2)*metric[0][0] - metric[1][1] - metric[2][2])
		metric[2][1] = metric[1][2]

		if metricInvalid() {
			log.Println("Metric invalid:")
			log.Println(metric)
			g22()
		}
	}
}

func g22() {
	// Control2 is +green
	S2Control2.Color = S2Control3.Color.Add(pixel.RGB(0, coloroffset, 0))

	ProgressButton.OnEvent = func(e *types.Event) {
		log.Println((S2Slider.Color.R - S2ControlColor.Color.R) / coloroffset)
		metric[1][1] = math.Pow((S2Slider.Color.R-S2ControlColor.Color.R)/coloroffset, 2) // Might not work when comparing locations
		g33()
	}
}

func metricInvalid() bool {
	if metric[1][2]*metric[1][2] >= metric[1][1]*metric[2][2] || metric[0][2]*metric[0][2] >= metric[0][0]*metric[2][2] || metric[0][1]*metric[0][1] >= metric[0][0]*metric[1][1] {
		return true
	}
	return false
}
