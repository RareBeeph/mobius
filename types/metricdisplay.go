package types

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MetricDisplay struct {
	CurveDisplay

	CenterCol   pixel.RGBA
	Lengths     [3][3]float64
	ColorOffset float64
}

func (d *MetricDisplay) Draw(window *pixelgl.Window) {
	d.GuardSurface()

	max := float64(0)
	for i := range [3]any{} {
		if d.Lengths[i][i] > max {
			max = d.Lengths[i][i]
		}
	}

	vertexCols := []pixel.RGBA{d.CenterCol,
		d.CenterCol.Add(pixel.RGB(2*d.ColorOffset*d.Lengths[0][0], 0, 0)),
		d.CenterCol.Add(pixel.RGB(0, 2*d.ColorOffset*d.Lengths[1][1], 0)),
		d.CenterCol.Add(pixel.RGB(0, 0, 2*d.ColorOffset*d.Lengths[2][2]))}

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
				p, depth := d.ProjectParallel(c.Sub(d.CenterCol).Scaled(1 / (2 * d.ColorOffset)))
				p = p.Scaled(1 / max)
				depth /= max
				d.surface.Push(d.Center.Add(p))

				d.surface.Circle(500/(200-depth), 0)
				axialDistance += sampleoffset
			}
		}
	}

	copy(d.BasisMatrix[:], temp[:])

	d.surface.Draw(window)
}

func (d *MetricDisplay) lengthsToMap() (out [3][3]float64) {
	var angles [3]float64

	// TODO: rework this
	// Copied from second_tc.go; terrible practice but I'll fix it later
	angles[0] = math.Acos((math.Pow(d.Lengths[0][1], 2) - math.Pow(d.Lengths[0][0], 2) - math.Pow(d.Lengths[1][1], 2)) / (-2 * d.Lengths[0][0] * d.Lengths[1][1]))
	angles[1] = math.Acos((math.Pow(d.Lengths[0][2], 2) - math.Pow(d.Lengths[0][0], 2) - math.Pow(d.Lengths[2][2], 2)) / (-2 * d.Lengths[0][0] * d.Lengths[2][2]))
	angles[2] = math.Acos((math.Pow(d.Lengths[1][2], 2) - math.Pow(d.Lengths[1][1], 2) - math.Pow(d.Lengths[2][2], 2)) / (-2 * d.Lengths[1][1] * d.Lengths[2][2]))

	for i := range angles {
		if math.IsNaN(angles[i]) {
			angles[i] = math.Pi / 2 // Safeguard to make sin and cos give reasonable values
		}
	}

	// Arbitrarily assert that the red direction of the map is in the red direction of space
	// out[0][0] = d.Lengths[0][0]
	out[0][0] = 1

	// Arbitrarily assert that the green direction of the map is in the red-green plane
	// out[1][0] = d.Lengths[1][1] * math.Cos(angles[0])
	// out[1][1] = d.Lengths[1][1] * math.Sin(angles[0])
	out[0][1] = math.Cos(angles[0])
	out[1][1] = math.Sin(angles[0])

	// Calculate some angle related stuff
	phi := math.Acos((math.Cos(angles[2]) - math.Cos(angles[1])*math.Cos(angles[0])) / (math.Sin(angles[1]) * math.Sin(angles[0])))
	if math.IsNaN(phi) {
		phi = 0 // Safeguard
	}

	// Constrain the blue direction of the map relative to those, based on the given lengths
	// out[2][0] = d.Lengths[2][2] * math.Cos(angles[1])
	// out[2][1] = d.Lengths[2][2] * math.Sin(angles[1]) * math.Cos(phi)
	// out[2][2] = d.Lengths[2][2] * math.Sin(angles[1]) * math.Sin(phi)
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
