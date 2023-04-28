package model

import (
	"gorm.io/gorm"
)

type Midpoint struct {
	gorm.Model
	Startpoint Color `gorm:"foreignKey:StartpointID"`
	Midpoint   Color `gorm:"foreignKey:MidpointID"`
	Endpoint   Color `gorm:"foreignKey:EndpointID"`

	StartpointID int
	MidpointID   int
	EndpointID   int
}

var AllModels []interface{}

func init() {
	AllModels = []interface{}{
		Midpoint{},
		Color{},
	}
}
