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

		// debug
		types.SlopesOfMetric(metric)
	} else if i > 2 || j > 2 {
		log.Println("Excessive arguments given to g(i,j); this should never happen")
	}

	toAdd := []float64{0, 0, 0}
	toAdd[i] = coloroffset
	if j != i {
		toAdd[j] = -coloroffset
	}

	tempCC3 := S2Control3.(*types.ColoredRect).Color
	if i == 0 && j == 0 {
		// Special case: when testing red's length, compare against a constant pair of grays
		// This should make metrics for different locations in color space directly comparable
		S2Control3.(*types.ColoredRect).Color = pixel.RGB(0.5, 0.5, 0.5)
		toAdd[1] = coloroffset
		toAdd[2] = coloroffset

		S2FpsCounter.(*types.FpsCounter).StepCount = 0
	}

	S2Control2.(*types.ColoredRect).Color = S2Control3.(*types.ColoredRect).Color.Add(pixel.RGB(toAdd[0], toAdd[1], toAdd[2]))

	ProgressButton.(*types.Button).OnEvent = func(e *types.Event) {
		S2FpsCounter.(*types.FpsCounter).StepCount++

		lengths[i][j] *= float64(timesLooped)
		if i == 0 && j == 0 {
			// length of red in terms of gray (Gy/R), as reciprocal of length of gray in terms of red (R/Gy)
			// TODO: Handle when red distance is 0. This will always snap lengths[0][0] to Inf.
			lengths[i][j] += math.Abs(coloroffset / (S2Slider.(*types.Slider).Color.R - S2ControlColor.(*types.ColoredRect).Color.R))

			// revert control color. probably doable just by setting to the first control color
			S2Control3.(*types.ColoredRect).Color = tempCC3
		} else {
			// length of X in terms of gray (Gy/X), as length of X in terms of red times length of red in terms of gray ((R/X)*(Gy/R))
			lengths[i][j] += math.Abs((S2Slider.(*types.Slider).Color.R - S2ControlColor.(*types.ColoredRect).Color.R) * lengths[0][0] / coloroffset)
		}
		lengths[i][j] /= float64(timesLooped + 1)

		copy(MetricGraph.(*types.MetricDisplay).Lengths[:], lengths[:])

		angles = calculateAngles()

		log.Println(lengths[i][j])
		if i == j {
			// Dot product of a vector with itself is its length squared
			metric[i][j] = math.Pow(lengths[i][j], 2)

			// On progression button pressed, "recurse", testing the next diagonal entry i+1==j+1
			// (Overflow will be corrected)
			measureMetric(i+1, j+1, timesLooped)
		} else {
			// Dot product of a vector with another, as a function of their lengths (rearranged law of cosines)
			metric[i][j] = -0.5 * (math.Pow(lengths[i][j], 2) - metric[i][i] - metric[j][j])
			metric[j][i] = metric[i][j]

			// On progression button pressed, "recurse", testing the next off-diagonal entry
			// (Overflow will be corrected)
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
	// Consider using the IsNaN formulation used in metricdisplay.go; consider consolidating into one implementation
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

func calculateAngles() (out [3]float64) {
	for i, j := range [][]int{{0, 1}, {0, 2}, {1, 2}} {
		// TODO: find a way to make this not go off the screen
		// Rearranged law of cosines to solve for angle
		out[i] = math.Acos((math.Pow(lengths[j[0]][j[1]], 2) - math.Pow(lengths[j[0]][j[0]], 2) - math.Pow(lengths[j[1]][j[1]], 2)) / (-2 * lengths[j[0]][j[0]] * lengths[j[1]][j[1]]))
	}
	return out
}

var lengths [3][3]float64
var angles [3]float64
