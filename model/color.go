package model

import (
	"github.com/faiface/pixel"
	"gorm.io/gorm"
)

type Color struct {
	gorm.Model
	R float64 `gorm:"index:idx_together,unique"`
	G float64 `gorm:"index:idx_together,unique"`
	B float64 `gorm:"index:idx_together,unique"`
}

func RgbaToColor(inputRGBA pixel.RGBA) Color {
	return Color{
		R: inputRGBA.R,
		G: inputRGBA.G,
		B: inputRGBA.B,
	}
}
