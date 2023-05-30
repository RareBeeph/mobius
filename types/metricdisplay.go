package types

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MetricDisplay struct {
	Entity

	Center      pixel.Vec
	Bounds      pixel.Rect
	BasisMatrix [3][3]float64

	CenterCol       pixel.RGBA
	Lengths         [3][3]float64
	ColorOffset     float64
	ThicknessFactor float64
	CenterDepth     float64
}

const sampleoffset = 0.01

var DefaultBasisMatrix = [3][3]float64{{-100 * math.Sqrt(0.5), 0, 100 * math.Sqrt(0.5)}, {-100 * math.Sqrt(0.16667), 100 * math.Sqrt(0.6667), -100 * math.Sqrt(0.16667)}, {100 * math.Sqrt(0.3333), 100 * math.Sqrt(0.3333), 100 * math.Sqrt(0.3333)}}

func (d *MetricDisplay) Draw(window *pixelgl.Window) {
	d.GuardSurface()

	max := float64(0)
	for i := range [3]any{} {
		if d.Lengths[i][i] > max {
			max = d.Lengths[i][i]
		}
	}

	// Use 3*d.ColorOffset to exaggerate the visual
	// TODO: set the exaggeration factor dynamically
	vertexCols := []pixel.RGBA{d.CenterCol,
		d.CenterCol.Add(pixel.RGB(3*d.ColorOffset*d.Lengths[0][0], 0, 0)),
		d.CenterCol.Add(pixel.RGB(0, 3*d.ColorOffset*d.Lengths[1][1], 0)),
		d.CenterCol.Add(pixel.RGB(0, 0, 3*d.ColorOffset*d.Lengths[2][2])),
	}

	// Temporarily store basis matrix (janky hack)
	var temp [3][3]float64
	copy(temp[:], d.BasisMatrix[:])

	// Transformation determined by lengths matrix
	d.BasisMatrix = dirtyMatrixMultiply(d.BasisMatrix, d.lengthsToMap())

	for i, col := range vertexCols {
		for j, col2 := range vertexCols {
			// Restrict rendering to only the 6 valid and distinct unordered pairs (edges)
			if j <= i {
				continue
			}

			// TODO: Render failed triangle inequalities differently or not at all

			axialDistance := float64(0)

			// TODO: Add a depth buffer or something here
			for axialDistance <= 1 {
				// Color of the point determined by linear interpolation between endpoint
				c := col.Add(col2.Sub(col).Scaled(axialDistance))
				d.surface.Color = c

				// Position of the point determined by basis matrix
				p, depth := d.ProjectParallel(c.Sub(d.CenterCol).Scaled(1 / (3 * d.ColorOffset)))
				p = p.Scaled(1 / max)
				depth /= max
				d.surface.Push(d.Center.Add(p))

				// TODO: set camera distance dynamically
				// Note: behaves poorly when center depth and point depth are too similar.
				d.surface.Circle(d.CenterDepth*d.ThicknessFactor/(d.CenterDepth-depth), 0)
				axialDistance += sampleoffset
			}
		}
	}

	copy(d.BasisMatrix[:], temp[:])

	d.surface.Draw(window)
}

func (d *MetricDisplay) Handles(delta *Event) bool {
	if delta.EventType != Drag {
		return false
	}
	if !d.Contains(delta.InitialPos) && !delta.Contains(pixelgl.KeyC) {
		return false
	}
	if !(delta.Contains(pixelgl.MouseButton1) || delta.Contains(pixelgl.MouseButton2) || delta.Contains(pixelgl.KeyC)) {
		return false
	}
	return true
}

func (d *MetricDisplay) Handle(delta *Event) {
	if delta.Contains(pixelgl.KeyC) {
		copy(d.BasisMatrix[:], DefaultBasisMatrix[:])
		return
	}
	if delta.Contains(pixelgl.MouseButton1) && delta.Contains(pixelgl.MouseButton2) {
		d.Speen(delta)
		return
	}

	rotatedMatrix := [3][3]float64{}
	copy(rotatedMatrix[:], d.BasisMatrix[:])

	var rotPhase float64
	if delta.MouseVel.X != 0 {
		rotPhase = math.Atan(delta.MouseVel.Y / delta.MouseVel.X)
	} else if delta.MouseVel.Y != 0 {
		rotPhase = math.Pi / 2 * delta.MouseVel.Y / math.Abs(delta.MouseVel.Y)
	} else {
		return
	}
	if delta.Contains(pixelgl.MouseButton2) {
		rotPhase = 0
	}

	rotMagnitude := delta.MouseVel.Len() / 100

	var rotPhaseSign float64 = 0
	if delta.MouseVel.X != 0 {
		rotPhaseSign = delta.MouseVel.X / math.Abs(delta.MouseVel.X)
	}

	rotVec := [3]float64{}

	rotVec[0] = math.Sin(rotMagnitude) * math.Cos(rotPhase) * rotPhaseSign

	if rotPhaseSign != 0 {
		rotVec[1] = math.Sin(rotMagnitude) * math.Sin(rotPhase) * rotPhaseSign
	} else {
		rotVec[1] = math.Sin(rotMagnitude) * delta.MouseVel.Y / math.Abs(delta.MouseVel.Y)
		if rotPhase == 0 {
			rotVec[1] = 0
		}
	}

	rotVec[2] = math.Sqrt(1 - rotVec[0]*rotVec[0] - rotVec[1]*rotVec[1])

	// Rotate with rotor based on 0,0,1 wedge rotVec
	for i := range [3]bool{} {
		rotatedMatrix[0][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[0][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[0][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[0][i] - 2*rotVec[0]*rotVec[1]*d.BasisMatrix[1][i] + 2*rotVec[0]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[1][i] = rotVec[0]*rotVec[0]*d.BasisMatrix[1][i] - rotVec[1]*rotVec[1]*d.BasisMatrix[1][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[1][i] - 2*rotVec[0]*rotVec[1]*d.BasisMatrix[0][i] + 2*rotVec[1]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[2][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[2][i] - rotVec[1]*rotVec[1]*d.BasisMatrix[2][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[2][i] - 2*rotVec[0]*rotVec[2]*d.BasisMatrix[0][i] - 2*rotVec[1]*rotVec[2]*d.BasisMatrix[1][i]
	}

	d.BasisMatrix = rotatedMatrix
}

func (d *MetricDisplay) Speen(delta *Event) {
	rotatedMatrix := [3][3]float64{}
	copy(rotatedMatrix[:], d.BasisMatrix[:])

	rotVec := [3]float64{}

	rotVec[0] = math.Sin(delta.MouseVel.X / 100)
	rotVec[1] = math.Cos(delta.MouseVel.X / 100)
	rotVec[2] = 0

	// Rotate with rotor based on 0,1,0 wedge rotVec
	// TODO: generalize this as its own function so it's not nearly-repeated from Handle()
	for i := range [3]bool{} {
		rotatedMatrix[0][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[0][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[0][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[0][i] + 2*rotVec[0]*rotVec[1]*d.BasisMatrix[1][i] - 2*rotVec[0]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[1][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[1][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[1][i] - rotVec[2]*rotVec[2]*d.BasisMatrix[1][i] - 2*rotVec[0]*rotVec[1]*d.BasisMatrix[0][i] - 2*rotVec[1]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[2][i] = rotVec[0]*rotVec[0]*d.BasisMatrix[2][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[2][i] - rotVec[2]*rotVec[2]*d.BasisMatrix[2][i] - 2*rotVec[0]*rotVec[2]*d.BasisMatrix[0][i] + 2*rotVec[1]*rotVec[2]*d.BasisMatrix[1][i]
	}

	d.BasisMatrix = rotatedMatrix
}

func (d *MetricDisplay) ProjectParallel(col pixel.RGBA) (out pixel.Vec, depth float64) {
	// Standard matrix multiplication
	out.X = col.R*d.BasisMatrix[0][0] + col.G*d.BasisMatrix[0][1] + col.B*d.BasisMatrix[0][2]
	out.Y = col.R*d.BasisMatrix[1][0] + col.G*d.BasisMatrix[1][1] + col.B*d.BasisMatrix[1][2]
	depth = col.R*d.BasisMatrix[2][0] + col.G*d.BasisMatrix[2][1] + col.B*d.BasisMatrix[2][2]
	return out, depth
}

func (d *MetricDisplay) Contains(point pixel.Vec) (out bool) {
	return (point.X >= d.Bounds.Min.X &&
		point.X < d.Bounds.Max.X &&
		point.Y >= d.Bounds.Min.Y &&
		point.Y < d.Bounds.Max.Y)
}

func (d *MetricDisplay) lengthsToMap() (out [3][3]float64) {
	var angles [3]float64

	// TODO: this is still just copypasta from second_tc.go
	for i, j := range [][]int{{0, 1}, {0, 2}, {1, 2}} {
		angles[i] = math.Acos((math.Pow(d.Lengths[j[0]][j[1]], 2) - math.Pow(d.Lengths[j[0]][j[0]], 2) - math.Pow(d.Lengths[j[1]][j[1]], 2)) / (-2 * d.Lengths[j[0]][j[0]] * d.Lengths[j[1]][j[1]]))
	}

	for i := range angles {
		if math.IsNaN(angles[i]) {
			angles[i] = math.Pi / 2 // Safeguard to make sin and cos give reasonable values
		}
	}

	// Arbitrarily assert that the red direction of the map is in the red direction of space
	out[0][0] = 1

	// Arbitrarily assert that the green direction of the map is in the red-green plane
	out[0][1] = math.Cos(angles[0])
	out[1][1] = math.Sin(angles[0])

	// Calculate phase of blue based on the dot product of green and blue
	phi := math.Acos((math.Cos(angles[2]) - math.Cos(angles[1])*math.Cos(angles[0])) / (math.Sin(angles[1]) * math.Sin(angles[0])))
	// If that phase is NaN, the triangle inequality failed
	if math.IsNaN(phi) {
		phi = 0 // Safeguard
	}

	// Constrain the blue direction of the map relative to those, based on the given lengths
	out[0][2] = math.Cos(angles[1])
	out[1][2] = math.Sin(angles[1]) * math.Cos(phi)
	out[2][2] = math.Sin(angles[1]) * math.Sin(phi)

	return out
}

func dirtyMatrixMultiply(m1 [3][3]float64, m2 [3][3]float64) (out [3][3]float64) {
	// Wrote this so I didn't have to figure out using a library
	for i := range [3]any{} {
		for j := range [3]any{} {
			for k := range [3]any{} {
				out[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}
	return out
}

/*
type point struct {
	col   pixel.RGBA
	pos   pixel.Vec
	depth float64
}

func PointSort(p []point) []point {
	sorted := make([]point, len(p))
	copy(sorted, p)

	// Use our own compare function to sort by depth ascending
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].depth < sorted[j].depth
	})

	return sorted
}
*/
