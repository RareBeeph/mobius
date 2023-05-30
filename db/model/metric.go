package model

import (
	"github.com/faiface/pixel"
	"gorm.io/gorm"
)

type Metric struct {
	gorm.Model
	pixel.RGBA

	RedSquared, GreenSquared, BlueSquared, RedDotGreen, RedDotBlue, GreenDotBlue float64
}

func NewMetricFromArray(a [3][3]float64, c pixel.RGBA) (out Metric) {
	out.RedSquared = a[0][0]
	out.GreenSquared = a[1][1]
	out.BlueSquared = a[2][2]
	out.RedDotGreen = a[0][1]
	out.RedDotBlue = a[0][2]
	out.GreenDotBlue = a[1][2]

	out.RGBA = c

	return out
}

var AllModels = []interface{}{
	Metric{},
}
