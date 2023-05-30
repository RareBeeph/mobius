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

	// Unique index might be the way to go here, but I couldn't make it work. See scene_two.go MetricSaveButton declaration
	ControlR float64 // `gorm:"index:yomama,unique"`
	ControlG float64 // `gorm:"index:yomama,unique"`
	ControlB float64 // `gorm:"index:yomama,unique"`
}

func NewMetricFromArray(a [3][3]float64, c Color) (out Metric) {
	out.RedSquared = a[0][0]
	out.GreenSquared = a[1][1]
	out.BlueSquared = a[2][2]
	out.RedDotGreen = a[0][1]
	out.RedDotBlue = a[0][2]
	out.GreenDotBlue = a[1][2]

	out.ControlR = c.R
	out.ControlG = c.G
	out.ControlB = c.B

	return out
}
