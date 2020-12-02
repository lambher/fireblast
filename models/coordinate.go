package models

import "github.com/faiface/pixel"

type Coordinate struct {
	Position    pixel.Vec
	Translation pixel.Vec
}

func (c Coordinate) GetPos() pixel.Vec {
	return c.Position.Add(c.Translation)
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
