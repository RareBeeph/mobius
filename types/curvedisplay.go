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
	BasisMatrix [3][3]float64
	Bounds      pixel.Rect
}

type point struct {
	col   pixel.RGBA
	pos   pixel.Vec
	depth float64
}

const sampleoffset = 0.01

var DefaultBasisMatrix = [3][3]float64{{-70, 0, 70}, {-50, 100, -50}, {50, 50, 50}}

func (d *CurveDisplay) Draw(window *pixelgl.Window) {
	d.GuardSurface()

	var pointList []point

	samplingProgress := float64(0)
	for samplingProgress <= 1 {
		pointList = append(pointList, point{col: d.Curve.EvenLagrangeInterp(samplingProgress)})
		pos, depth := d.ProjectParallel(pointList[len(pointList)-1].col)
		pointList[len(pointList)-1].pos = pos
		pointList[len(pointList)-1].depth = depth
		samplingProgress += sampleoffset
	}

	pointList = PointSort(pointList)

	// Gridlines. Will render behind curve regardless of depth; fix that later
	for i, col := range []pixel.RGBA{pixel.RGB(1, 0, 0), pixel.RGB(0, 1, 0), pixel.RGB(0, 0, 1)} {
		axialDistance := float64(0)
		for axialDistance <= 1 {
			d.surface.Color = col.Scaled(axialDistance)
			d.surface.Push(d.Center.Add(pixel.V(d.BasisMatrix[0][i], d.BasisMatrix[1][i]).Scaled(axialDistance)))
			d.surface.Circle(300/(200-d.BasisMatrix[2][i]*axialDistance), 0)
			axialDistance += sampleoffset
		}
	}

	for _, poi := range pointList {
		d.surface.Color = poi.col
		d.surface.Push(d.Center.Add(poi.pos))
		d.surface.Circle(500/(200-poi.depth), 0)
	}

	d.surface.Draw(window)
}

func (d *CurveDisplay) Receive(delta *Event) {
	// Duplicate, but without this here it'd run the generic Entity Handles() which is always false
	// Doesn't use the entity's new mutex stuff
	if d.Handles(delta) {
		d.Handle(delta)
	}
}

func (d *CurveDisplay) Handles(delta *Event) bool {
	if !d.Contains(delta.InitialPos) && !delta.Contains(pixelgl.KeyC) {
		return false
	}
	if !(delta.Contains(pixelgl.MouseButton1) || delta.Contains(pixelgl.MouseButton2) || delta.Contains(pixelgl.KeyC)) {
		return false
	}
	return true
}

func (d *CurveDisplay) Handle(delta *Event) {
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
	if delta.MousePos.X != 0 {
		rotPhase = math.Atan(delta.MousePos.Y / delta.MousePos.X)
	} else if delta.MousePos.Y != 0 {
		rotPhase = math.Pi / 2 * delta.MousePos.Y / math.Abs(delta.MousePos.Y)
	} else {
		return
	}
	if delta.Contains(pixelgl.MouseButton2) {
		rotPhase = 0
	}

	rotMagnitude := delta.MousePos.Len() / 100

	var rotPhaseSign float64 = 0
	if delta.MousePos.X != 0 {
		rotPhaseSign = delta.MousePos.X / math.Abs(delta.MousePos.X)
	}

	rotVec := [3]float64{}

	rotVec[0] = math.Sin(rotMagnitude) * math.Cos(rotPhase) * rotPhaseSign

	if rotPhaseSign != 0 {
		rotVec[1] = math.Sin(rotMagnitude) * math.Sin(rotPhase) * rotPhaseSign
	} else {
		rotVec[1] = math.Sin(rotMagnitude) * delta.MousePos.Y / math.Abs(delta.MousePos.Y)
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

func (d *CurveDisplay) Speen(delta *Event) {
	rotatedMatrix := [3][3]float64{}
	copy(rotatedMatrix[:], d.BasisMatrix[:])

	rotVec := [3]float64{}

	rotVec[0] = math.Sin(delta.MousePos.X / 100)
	rotVec[1] = math.Cos(delta.MousePos.X / 100)
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

func PointSort(p []point) []point {
	sorted := make([]point, len(p))
	copy(sorted, p)

	// Use our own compare function to sort by depth ascending
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].depth < sorted[j].depth
	})

	return sorted
}

func (d *CurveDisplay) ProjectParallel(col pixel.RGBA) (out pixel.Vec, depth float64) {
	// Standard matrix multiplication
	out.X = col.R*d.BasisMatrix[0][0] + col.G*d.BasisMatrix[0][1] + col.B*d.BasisMatrix[0][2]
	out.Y = col.R*d.BasisMatrix[1][0] + col.G*d.BasisMatrix[1][1] + col.B*d.BasisMatrix[1][2]
	depth = col.R*d.BasisMatrix[2][0] + col.G*d.BasisMatrix[2][1] + col.B*d.BasisMatrix[2][2]
	return out, depth
}

func (d *CurveDisplay) Contains(point pixel.Vec) (out bool) {
	return (point.X >= d.Bounds.Min.X &&
		point.X < d.Bounds.Max.X &&
		point.Y >= d.Bounds.Min.Y &&
		point.Y < d.Bounds.Max.Y)
}
