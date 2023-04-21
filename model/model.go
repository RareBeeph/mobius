package model

import (
	"github.com/faiface/pixel"
	"gorm.io/gorm"
)

type Midpoint struct {
	gorm.Model
	startpoint pixel.RGBA
	endpoint   pixel.RGBA
	midpoint   pixel.RGBA
}
