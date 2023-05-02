package types

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
	if d.surface == nil {
		d.surface = imdraw.New(nil)
	}

	d.surface.Clear()

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

	for _, poi := range pointlist {
		d.surface.Color = poi.col
		d.surface.Push(poi.pos)
		d.surface.Circle(10, 0)
	}

	d.surface.Draw(window)
}

func InsertionSort(p []point) []point {
	// Insertion sort is the first one that came to mind.
	// I know it's slow conceptually and my implementation probably sucks on top, but I don't know to sort these by depth with, say, sort package

	var mindepth float64
	var mindex int

	for i, poi := range p {
		mindepth = math.Inf(1)

		for j, poin := range p[i:] {
			if poin.depth < mindepth {
				mindepth = poin.depth
				mindex = j
			}
		}

		p[i] = p[mindex]
		p[mindex] = poi
	}
	return p // Temp
}

func (d *CurveDisplay) ProjectParallel(col pixel.RGBA) (out pixel.Vec, depth float64) {
	out.X = col.R*d.BasisMatrix[0] + col.G*d.BasisMatrix[1] + col.B*d.BasisMatrix[2]
	out.Y = col.R*d.BasisMatrix[3] + col.G*d.BasisMatrix[4] + col.B*d.BasisMatrix[5]
	depth = col.R*d.BasisMatrix[6] + col.G*d.BasisMatrix[7] + col.B*d.BasisMatrix[8]
	return out, depth
}
