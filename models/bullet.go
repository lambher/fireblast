package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/conf"
)

type Bullet struct {
	Conf *conf.Conf

	destroyed bool
	Color     pixel.RGBA
	Position  *Coordinate
	Inertie   pixel.Vec
}

func NewBullet(Conf *conf.Conf) *Bullet {
	bullet := Bullet{
		Conf:  Conf,
		Color: pixel.RGB(1, 1, 1),
	}
	return &bullet
}

func (b Bullet) Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	imd.Color = b.Color
	imd.Push(b.Position.Position.Add(b.Position.Translation))
	imd.Circle(1, 0)
	imd.Draw(win)
}

func (b *Bullet) applyTranslation() {
	b.Translate(b.Inertie)
}

func (b *Bullet) Translate(vec pixel.Vec) {
	b.Position.Translation = b.Position.Translation.Add(vec)
}
func (b *Bullet) Update() {
	b.applyTranslation()
	if b.Position.Position.Add(b.Position.Translation).X > b.Conf.MaxX || b.Position.Position.Add(b.Position.Translation).Y > b.Conf.MaxY {
		b.destroyed = true
	}
}

func (p Bullet) Collides(element Element) bool {
	return false
}

func (b *Bullet) Collision(elements []Element) {
}

func (b Bullet) IsDestroyed() bool {
	return b.destroyed
}
