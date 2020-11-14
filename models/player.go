package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/conf"
)

type Player struct {
	Conf *conf.Conf

	Name     string
	HP       int
	Position *Coordinate
	Color    pixel.RGBA

	Shape *Triangle

	gazActive bool
}

type PlayerConf struct {
	Name     string
	Position pixel.Vec
	Color    pixel.RGBA
	Conf     *conf.Conf
}

func NewPlayer(conf PlayerConf) *Player {
	player := Player{
		Conf:     conf.Conf,
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
	if p.gazActive {
		p.Shape.Gaz()
	}
	p.Shape.Update()
}

func (p *Player) Shoot() *Bullet {
	bullet := NewBullet(p.Conf)
	bullet.Position = NewCoordinate(p.Shape.A.Translation)
	bullet.Inertie = p.Shape.Direction.Scaled(5).Add(p.Shape.Inertie)
	return bullet
}

func (p *Player) Gaz(active bool) {
	p.gazActive = active
}

func (p *Player) Rotate(teta float64) {
	p.Shape.RotateDirection(teta)
}

func (p Player) IsDestroyed() bool {
	return false
}
