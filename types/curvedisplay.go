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
	Bounds      pixel.Rect
}

type point struct {
	col   pixel.RGBA
	pos   pixel.Vec
	depth float64
}

const GRAIN = 0.01

var DefaultBasisMatrix = [9]float64{-70, 0, 70, -50, 100, -50, 50, 50, 50}

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

	// Gridlines. Will render behind curve regardless of depth; fix that later
	for i, col := range []pixel.RGBA{pixel.RGB(1, 0, 0), pixel.RGB(0, 1, 0), pixel.RGB(0, 0, 1)} {
		t := float64(0)
		for t <= 1 {
			d.surface.Color = col.Scaled(t)
			d.surface.Push(d.Center.Add(pixel.V(d.BasisMatrix[i], d.BasisMatrix[3+i]).Scaled(t)))
			d.surface.Circle(300/(200-d.BasisMatrix[6+i]*t), 0)
			t += GRAIN
		}
	}

	for _, poi := range pointlist {
		d.surface.Color = poi.col
		d.surface.Push(d.Center.Add(poi.pos))
		d.surface.Circle(500/(200-poi.depth), 0)
	}

	d.surface.Draw(window)
}

func (d *CurveDisplay) Handle(delta Event) {
	if !d.Contains(delta.InitialPos) {
		return
	}
	if !(delta.Contains(pixelgl.MouseButton1) || delta.Contains(pixelgl.MouseButton2) || delta.Contains(pixelgl.KeyC)) {
		return
	}
	if delta.Contains(pixelgl.KeyC) {
		copy(d.BasisMatrix[:], DefaultBasisMatrix[:])
		return
	}
	if delta.Contains(pixelgl.MouseButton1) && delta.Contains(pixelgl.MouseButton2) {
		d.Speen(delta)
		return
	}

	rotatedMatrix := [9]float64{}
	copy(rotatedMatrix[:], d.BasisMatrix[:])

	var phi float64
	if delta.MousePos.X != 0 {
		phi = math.Atan(delta.MousePos.Y / delta.MousePos.X)
	} else if delta.MousePos.Y != 0 {
		phi = math.Pi / 2 * delta.MousePos.Y / math.Abs(delta.MousePos.Y)
	} else {
		return
	}
	if delta.Contains(pixelgl.MouseButton2) {
		phi = 0
	}

	theta := delta.MousePos.Len() / 100

	var sign float64 = 0
	if delta.MousePos.X != 0 {
		sign = delta.MousePos.X / math.Abs(delta.MousePos.X)
	}

	rotvec := [3]float64{}

	rotvec[0] = math.Sin(theta) * math.Cos(phi) * sign

	if sign != 0 {
		rotvec[1] = math.Sin(theta) * math.Sin(phi) * sign
	} else {
		rotvec[1] = math.Sin(theta) * delta.MousePos.Y / math.Abs(delta.MousePos.Y)
		if phi == 0 {
			rotvec[1] = 0
		}
	}

	rotvec[2] = math.Sqrt(1 - rotvec[0]*rotvec[0] - rotvec[1]*rotvec[1])

	// Rotate with rotor based on 0,0,1 wedge rotvec
	for i := range [3]bool{} {
		rotatedMatrix[i] = -rotvec[0]*rotvec[0]*d.BasisMatrix[i] + rotvec[1]*rotvec[1]*d.BasisMatrix[i] + rotvec[2]*rotvec[2]*d.BasisMatrix[i] - 2*rotvec[0]*rotvec[1]*d.BasisMatrix[3+i] + 2*rotvec[0]*rotvec[2]*d.BasisMatrix[6+i]
		rotatedMatrix[3+i] = rotvec[0]*rotvec[0]*d.BasisMatrix[3+i] - rotvec[1]*rotvec[1]*d.BasisMatrix[3+i] + rotvec[2]*rotvec[2]*d.BasisMatrix[3+i] - 2*rotvec[0]*rotvec[1]*d.BasisMatrix[i] + 2*rotvec[1]*rotvec[2]*d.BasisMatrix[6+i]
		rotatedMatrix[6+i] = -rotvec[0]*rotvec[0]*d.BasisMatrix[6+i] - rotvec[1]*rotvec[1]*d.BasisMatrix[6+i] + rotvec[2]*rotvec[2]*d.BasisMatrix[6+i] - 2*rotvec[0]*rotvec[2]*d.BasisMatrix[i] - 2*rotvec[1]*rotvec[2]*d.BasisMatrix[3+i]
	}

	d.BasisMatrix = rotatedMatrix
}

func (d *CurveDisplay) Speen(delta Event) {
	rotatedMatrix := [9]float64{}
	copy(rotatedMatrix[:], d.BasisMatrix[:])

	rotvec := [3]float64{}

	rotvec[1] = math.Cos(delta.MousePos.X / 100)
	rotvec[0] = math.Sin(delta.MousePos.X / 100)
	rotvec[2] = 0

	// Rotate with rotor based on 0,1,0 wedge rotvec
	for i := range [3]bool{} {
		rotatedMatrix[i] = -rotvec[0]*rotvec[0]*d.BasisMatrix[i] + rotvec[1]*rotvec[1]*d.BasisMatrix[i] + rotvec[2]*rotvec[2]*d.BasisMatrix[i] + 2*rotvec[0]*rotvec[1]*d.BasisMatrix[3+i] - 2*rotvec[0]*rotvec[2]*d.BasisMatrix[6+i]
		rotatedMatrix[3+i] = -rotvec[0]*rotvec[0]*d.BasisMatrix[3+i] + rotvec[1]*rotvec[1]*d.BasisMatrix[3+i] - rotvec[2]*rotvec[2]*d.BasisMatrix[3+i] - 2*rotvec[0]*rotvec[1]*d.BasisMatrix[i] - 2*rotvec[1]*rotvec[2]*d.BasisMatrix[6+i]
		rotatedMatrix[6+i] = rotvec[0]*rotvec[0]*d.BasisMatrix[6+i] + rotvec[1]*rotvec[1]*d.BasisMatrix[6+i] - rotvec[2]*rotvec[2]*d.BasisMatrix[6+i] - 2*rotvec[0]*rotvec[2]*d.BasisMatrix[i] + 2*rotvec[1]*rotvec[2]*d.BasisMatrix[3+i]
	}

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

func (d *CurveDisplay) Contains(point pixel.Vec) (out bool) {
	return (point.X >= d.Bounds.Min.X &&
		point.X < d.Bounds.Max.X &&
		point.Y >= d.Bounds.Min.Y &&
		point.Y < d.Bounds.Max.Y)
}
