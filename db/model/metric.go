package model

import "gorm.io/gorm"

type Metric struct {
	gorm.Model

	RedSquared   float64
	GreenSquared float64
	BlueSquared  float64
	RedDotGreen  float64
	RedDotBlue   float64
	GreenDotBlue float64

	ControlColor Color `gorm:"foreignKey:ControlID"`
	ControlID    int
}

func NewMetricFromArray(a [3][3]float64, c Color) (out Metric) {
	out.RedSquared = a[0][0]
	out.GreenSquared = a[1][1]
	out.BlueSquared = a[2][2]
	out.RedDotGreen = a[0][1]
	out.RedDotBlue = a[0][2]
	out.GreenDotBlue = a[1][2]

	out.ControlColor = c

	return out
}
