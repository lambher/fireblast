package models

import "github.com/faiface/pixel"

type Coordinate struct {
	Position    pixel.Vec
	Translation pixel.Vec
}

func NewCoordinate(vec pixel.Vec) *Coordinate {
	c := Coordinate{
		Position: vec,
		Translation: pixel.Vec{
			X: 0,
			Y: 0,
		},
	}
	return &c
}
