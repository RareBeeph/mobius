package types

import (
	"math"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type CurveDisplay struct {
	Entity

	Center      pixel.Vec
	Curve       RgbCurve
	BasisMatrix [9]float64
}

type point struct {
	col   pixel.RGBA
	pos   pixel.Vec
	depth float64
}

const GRAIN = 0.01

func (d *CurveDisplay) Draw(window *pixelgl.Window) {
	d.GuardSurface()

	var pointlist []point

	t := float64(0)
	for t <= 1 {
		pointlist = append(pointlist, point{col: d.Curve.EvenLagrangeInterp(t)})
		a, b := d.ProjectParallel(pointlist[len(pointlist)-1].col) // Temp
		pointlist[len(pointlist)-1].pos = a
		pointlist[len(pointlist)-1].depth = b
		t += GRAIN
	}

	pointlist = InsertionSort(pointlist) // Temp

	for i, col := range []pixel.RGBA{pixel.RGB(1, 0, 0), pixel.RGB(0, 1, 0), pixel.RGB(0, 0, 1)} {
		d.surface.Color = pixel.RGB(0, 0, 0)
		d.surface.Push(d.Center)
		d.surface.Color = col
		d.surface.Push(d.Center.Add(pixel.V(d.BasisMatrix[i], d.BasisMatrix[3+i])))
		d.surface.Line(2)
	}

	for _, poi := range pointlist {
		d.surface.Color = poi.col
		d.surface.Push(d.Center.Add(poi.pos))
		d.surface.Circle(500/(200-poi.depth), 0)
	}

	d.surface.Draw(window)
}

func (d *CurveDisplay) Handle(delta Event) {
	rotatedMatrix := d.BasisMatrix

	theta := delta.MousePos.X / 100

	// Fully bodged, I didn't even do the actual matrix math I just kinda messed around until it looked good
	rotatedMatrix[0] = d.BasisMatrix[0]*math.Cos(theta) + d.BasisMatrix[6]*math.Sin(theta)
	rotatedMatrix[6] = -d.BasisMatrix[0]*math.Sin(theta) + d.BasisMatrix[6]*math.Cos(theta)

	rotatedMatrix[1] = d.BasisMatrix[1]*math.Cos(theta) + d.BasisMatrix[7]*math.Sin(theta)
	rotatedMatrix[7] = -d.BasisMatrix[1]*math.Sin(theta) + d.BasisMatrix[7]*math.Cos(theta)

	rotatedMatrix[2] = d.BasisMatrix[2]*math.Cos(theta) + d.BasisMatrix[8]*math.Sin(theta)
	rotatedMatrix[8] = -d.BasisMatrix[2]*math.Sin(theta) + d.BasisMatrix[8]*math.Cos(theta)

	d.BasisMatrix = rotatedMatrix
}

func InsertionSort(p []point) []point {
	// Create a copy of the slice so we're not modifying an input argument
	sorted := make([]point, len(p))
	copy(sorted, p)

	// Use our own compare function to sort by depth ascending
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].depth < sorted[j].depth
	})

	return sorted
}

func (d *CurveDisplay) ProjectParallel(col pixel.RGBA) (out pixel.Vec, depth float64) {
	out.X = col.R*d.BasisMatrix[0] + col.G*d.BasisMatrix[1] + col.B*d.BasisMatrix[2]
	out.Y = col.R*d.BasisMatrix[3] + col.G*d.BasisMatrix[4] + col.B*d.BasisMatrix[5]
	depth = col.R*d.BasisMatrix[6] + col.G*d.BasisMatrix[7] + col.B*d.BasisMatrix[8]
	return out, depth
}
