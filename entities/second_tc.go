package entities

import (
	"colorspacer/types"
	"log"
	"math"

	"github.com/faiface/pixel"
)

func measureMetric(i int, j int, timesLooped int) {
	if i < 0 || j < 0 || timesLooped < 0 {
		log.Println("Negative arguments given to measureMetric; this should never happen")
		return
	} else if j == 3 && i == 3 {
		log.Println("Switching to measuring diagonals")
		i = 0
		j = 1
	} else if i == 0 && j == 3 {
		log.Println("Switching to measuring +G-B")
		i = 1
		j = 2
	} else if i == 1 && j == 3 {
		if metricInvalid() {
			log.Println("Metric invalid:")
			log.Println(metric)
		}
		log.Println("Looping")
		i = 0
		j = 0
		timesLooped++
	} else if i > 2 || j > 2 {
		log.Println("Excessive arguments given to g(i,j); this should never happen")
	}

	toAdd := []float64{0, 0, 0}
	toAdd[i] = coloroffset
	if j != i {
		toAdd[j] = -coloroffset
	}
	S2Control2.Color = S2Control3.Color.Add(pixel.RGB(toAdd[0], toAdd[1], toAdd[2]))

	ProgressButton.OnEvent = func(e *types.Event) {
		lengths[i][j] *= float64(timesLooped)
		lengths[i][j] += math.Abs((S2Slider.Color.R - S2ControlColor.Color.R) / coloroffset /* * math.Sqrt(metric[0][0]) */)
		lengths[i][j] /= float64(timesLooped + 1)

		copy(MetricGraph.Lengths[:], lengths[:])

		calculateAngles()

		log.Println(lengths[i][j])
		if i == j {
			metric[i][j] = math.Pow(lengths[i][j], 2)

			measureMetric(i+1, j+1, timesLooped)
		} else {
			metric[i][j] = -0.5 * (math.Pow(lengths[i][j], 2) /* *metric[0][0] */ - metric[i][i] - metric[j][j])
			metric[j][i] = metric[i][j]

			measureMetric(i, j+1, timesLooped)
		}
	}
}

func metricInvalid() bool {
	// TODO: simplify these conditions

	// This check was previously done based on the triangle inequality with metric elements, but I missed a condition. This is easier.
	for i := range angles {
		if math.IsNaN(angles[i]) {
			log.Println("Metric invalid: NaN angles")
			log.Println(angles)
			return true
		}
	}

	var modifiedAngles [3]float64
	// TODO: make sure this is a valid formulation of the spherical triangle inequality for all angles greater than pi/2
	// I think it is but I'm not actually 100%
	for i := range angles {
		modifiedAngles[i] = (-math.Abs(angles[i]-math.Pi/2) + math.Pi/2)
	}
	if modifiedAngles[0]+modifiedAngles[1] < modifiedAngles[2] || modifiedAngles[0]+modifiedAngles[2] < modifiedAngles[1] || modifiedAngles[1]+modifiedAngles[2] < modifiedAngles[0] {
		log.Println("Metric invalid: Unsatisfiable angular constraints")
		log.Println(modifiedAngles)
		return true
	}

	return false
}

func calculateAngles() {
	// TODO: rework this
	angles[0] = math.Acos((math.Pow(lengths[0][1], 2) - math.Pow(lengths[0][0], 2) - math.Pow(lengths[1][1], 2)) / (-2 * lengths[0][0] * lengths[1][1]))
	angles[1] = math.Acos((math.Pow(lengths[0][2], 2) - math.Pow(lengths[0][0], 2) - math.Pow(lengths[2][2], 2)) / (-2 * lengths[0][0] * lengths[2][2]))
	angles[2] = math.Acos((math.Pow(lengths[1][2], 2) - math.Pow(lengths[1][1], 2) - math.Pow(lengths[2][2], 2)) / (-2 * lengths[1][1] * lengths[2][2]))
}

var lengths [3][3]float64
var angles [3]float64
