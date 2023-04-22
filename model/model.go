package model

import (
	"gorm.io/gorm"
)

type Midpoint struct {
	gorm.Model
	StartpointR float64
	StartpointG float64
	StartpointB float64
	EndpointR   float64
	EndpointG   float64
	EndpointB   float64
	MidpointR   float64
	MidpointG   float64
	MidpointB   float64
}
