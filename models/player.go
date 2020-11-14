package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Player struct {
	Name     string
	HP       int
	Position *Coordinate
	Color    pixel.RGBA

	Shape *Triangle
}

type PlayerConf struct {
	Name     string
	Position pixel.Vec
	Color    pixel.RGBA
}

func NewPlayer(conf PlayerConf) *Player {
	player := Player{
		Name:     conf.Name,
		HP:       100,
		Position: NewCoordinate(conf.Position),
		Color:    conf.Color,
	}

	player.Shape = NewTriangle(&player)

	return &player
}

func (p Player) Draw(win *pixelgl.Window) {
	p.Shape.Draw(win)
}

func (p *Player) Update() {
	p.Shape.Update()
}

func (p *Player) Rotate(teta float64) {
	p.Shape.RotateDirection(teta)
}
