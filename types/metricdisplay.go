package types

import (
	"math"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MobiusDisplay struct {
	Entity

	Center      pixel.Vec
	Bounds      pixel.Rect
	BasisMatrix [3][3]float64

	ThicknessFactor float64
	CenterDepth     float64
}

const sampleoffset = 0.01

var DefaultBasisMatrix = [3][3]float64{{-100 * math.Sqrt(0.5), 0, 100 * math.Sqrt(0.5)}, {-100 * math.Sqrt(0.16667), 100 * math.Sqrt(0.6667), -100 * math.Sqrt(0.16667)}, {100 * math.Sqrt(0.3333), 100 * math.Sqrt(0.3333), 100 * math.Sqrt(0.3333)}}

const size = 0.5

const PI = 3.1415926535897932384626

func getColor(theta float64, i float64) pixel.RGBA {
	// theta,i in half-revolutions
	hue := (theta + i/2) * PI

	// temp
	e1 := []float64{1 / math.Sqrt(2), -1 / math.Sqrt(2), 0}
	e2 := []float64{math.Sqrt(2./3.) / 2., math.Sqrt(2./3.) / 2., -math.Sqrt(2. / 3.)}

	rscale := 0.75
	gscale := 0.5
	bscale := 1.

	return pixel.RGBA{
		R: 0.5 + rscale*(e1[0]*math.Cos(hue)+e2[0]*math.Sin(hue))/2,
		G: 0.5 + gscale*(e1[1]*math.Cos(hue)+e2[1]*math.Sin(hue))/2,
		B: 0.5 + bscale*(e1[2]*math.Cos(hue)+e2[2]*math.Sin(hue))/2,
		A: 1,
	}
}

const offsetamount = 0.5

func getPosition(theta float64, i float64) (pixel.Vec, float64) {
	// theta in revolutions, i linear
	theta *= 2 * PI

	base := []float64{math.Cos(theta), math.Sin(theta), 0}

	phi := theta / 2

	offset := []float64{math.Cos(theta) * math.Cos(phi), math.Sin(theta) * math.Cos(phi), math.Sin(phi)}

	return pixel.Vec{
		X: size * (base[0] + i*offset[0]*offsetamount),
		Y: size * (base[1] + i*offset[1]*offsetamount),
	}, size * (base[2] + i*offset[2]*offsetamount)
}

func getPoint(theta float64, i float64) point {
	a, b := getPosition(theta, i)
	return point{col: getColor(theta, i), pos: a, depth: b}
}

func getPoint2(theta float64, i float64) point {
	a, b := getPosition(theta, i)
	return point{col: getColor(theta+1, -i), pos: a, depth: b}
}

var tgrain = 31.
var Tgrain = &tgrain

var igrain = 10.
var Igrain = &igrain

func (d *MobiusDisplay) Draw(window *pixelgl.Window) {
	d.GuardSurface()

	tgrain = math.Round(tgrain)
	igrain = math.Round(igrain)

	points := []point{}
	points2 := []point{}
	for i := 0; i < int(igrain); i++ {
		for theta := 0; theta < int(tgrain); theta++ {
			hueangle := float64(theta) / tgrain
			offset := -1 + 2*float64(i)/(igrain-1)

			p1 := getPoint(hueangle, offset)
			p1.pos, p1.depth = d.ProjectParallel(pixel.RGBA{R: p1.pos.X, G: p1.pos.Y, B: p1.depth})
			points = append(points, p1)

			p2 := getPoint2(hueangle, offset)
			p2.pos, p2.depth = d.ProjectParallel(pixel.RGBA{R: p2.pos.X, G: p2.pos.Y, B: p2.depth})
			p2.pos.X += 200
			points2 = append(points2, p2)
		}
	}

	points = PointSort(points)
	points2 = PointSort(points2)

	for _, p := range points {
		d.surface.Color = p.col
		d.surface.Push(d.Center.Add(p.pos))
		d.surface.Circle(d.CenterDepth*d.ThicknessFactor/(d.CenterDepth-p.depth), 0)
	}

	for _, p := range points2 {
		d.surface.Color = p.col
		d.surface.Push(d.Center.Add(p.pos))
		d.surface.Circle(d.CenterDepth*d.ThicknessFactor/(d.CenterDepth-p.depth), 0)
	}

	d.surface.Draw(window)
}

func (d *MobiusDisplay) Handles(delta *Event) bool {
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

func (d *MobiusDisplay) Handle(delta *Event) {
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

	// Rotate with rotor based on 0,0,1 * rotVec
	for i := range [3]bool{} {
		rotatedMatrix[0][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[0][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[0][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[0][i] - 2*rotVec[0]*rotVec[1]*d.BasisMatrix[1][i] + 2*rotVec[0]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[1][i] = rotVec[0]*rotVec[0]*d.BasisMatrix[1][i] - rotVec[1]*rotVec[1]*d.BasisMatrix[1][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[1][i] - 2*rotVec[0]*rotVec[1]*d.BasisMatrix[0][i] + 2*rotVec[1]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[2][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[2][i] - rotVec[1]*rotVec[1]*d.BasisMatrix[2][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[2][i] - 2*rotVec[0]*rotVec[2]*d.BasisMatrix[0][i] - 2*rotVec[1]*rotVec[2]*d.BasisMatrix[1][i]
	}

	d.BasisMatrix = rotatedMatrix
}

func (d *MobiusDisplay) Speen(delta *Event) {
	rotatedMatrix := [3][3]float64{}
	copy(rotatedMatrix[:], d.BasisMatrix[:])

	rotVec := [3]float64{}

	rotVec[0] = math.Sin(delta.MouseVel.X / 100)
	rotVec[1] = math.Cos(delta.MouseVel.X / 100)
	rotVec[2] = 0

	// Rotate with rotor based on 0,1,0 * rotVec
	// TODO: generalize this as its own function so it's not nearly-repeated from Handle()
	for i := range [3]bool{} {
		rotatedMatrix[0][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[0][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[0][i] + rotVec[2]*rotVec[2]*d.BasisMatrix[0][i] + 2*rotVec[0]*rotVec[1]*d.BasisMatrix[1][i] - 2*rotVec[0]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[1][i] = -rotVec[0]*rotVec[0]*d.BasisMatrix[1][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[1][i] - rotVec[2]*rotVec[2]*d.BasisMatrix[1][i] - 2*rotVec[0]*rotVec[1]*d.BasisMatrix[0][i] - 2*rotVec[1]*rotVec[2]*d.BasisMatrix[2][i]
		rotatedMatrix[2][i] = rotVec[0]*rotVec[0]*d.BasisMatrix[2][i] + rotVec[1]*rotVec[1]*d.BasisMatrix[2][i] - rotVec[2]*rotVec[2]*d.BasisMatrix[2][i] - 2*rotVec[0]*rotVec[2]*d.BasisMatrix[0][i] + 2*rotVec[1]*rotVec[2]*d.BasisMatrix[1][i]
	}

	d.BasisMatrix = rotatedMatrix
}

func (d *MobiusDisplay) ProjectParallel(col pixel.RGBA) (out pixel.Vec, depth float64) {
	// Standard matrix multiplication
	out.X = col.R*d.BasisMatrix[0][0] + col.G*d.BasisMatrix[0][1] + col.B*d.BasisMatrix[0][2]
	out.Y = col.R*d.BasisMatrix[1][0] + col.G*d.BasisMatrix[1][1] + col.B*d.BasisMatrix[1][2]
	depth = col.R*d.BasisMatrix[2][0] + col.G*d.BasisMatrix[2][1] + col.B*d.BasisMatrix[2][2]
	return out, depth
}

func (d *MobiusDisplay) Contains(point pixel.Vec) (out bool) {
	return (point.X >= d.Bounds.Min.X &&
		point.X < d.Bounds.Max.X &&
		point.Y >= d.Bounds.Min.Y &&
		point.Y < d.Bounds.Max.Y)
}

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
