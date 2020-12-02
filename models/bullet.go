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
	Player    *Player
	Velocity  pixel.Vec
}

func NewBullet(Conf *conf.Conf, player *Player) *Bullet {
	bullet := Bullet{
		Player: player,
		Conf:   Conf,
		Color:  pixel.RGB(1, 1, 1),
	}
	return &bullet
}

func (b Bullet) Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	imd.Color = b.Color
	imd.Push(b.Position.Position.Add(b.Position.Translation))
	imd.Circle(2, 0)
	imd.Draw(win)
	//
	//hitBox := b.GetHitBox()
	//
	//imd.Color = pixel.RGB(1, 0, 0)
	//imd.Push(b.Position.Position.Add(b.Position.Translation.Add(pixel.V(hitBox.Min.X, 0))))
	//imd.Color = pixel.RGB(1, 0, 0)
	//imd.Push(b.Position.Position.Add(b.Position.Translation.Add(pixel.V(0, hitBox.Min.Y))))
	//imd.Color = pixel.RGB(0, 1, 0)
	//imd.Push(b.Position.Position.Add(b.Position.Translation.Add(pixel.V(hitBox.Max.X, 0))))
	//imd.Color = pixel.RGB(0, 1, 0)
	//imd.Push(b.Position.Position.Add(b.Position.Translation.Add(pixel.V(0, hitBox.Max.Y))))
	////imd.Color = pixel.RGB(0, 0, 1)
	////imd.Push(pixel.V(500, 700))
	//imd.Polygon(0)
	//imd.Draw(win)
}

func (b *Bullet) applyTranslation() {
	b.Translate(b.Velocity)
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

func (p *Bullet) Hit(element Element) {

}

func (b Bullet) GetHitBox() pixel.Rect {
	return pixel.Rect{
		Min: b.Position.Position.Add(b.Position.Translation.Add(pixel.Vec{-2, -2})),
		Max: b.Position.Position.Add(b.Position.Translation.Add(pixel.Vec{2, 2})),
	}
}

func (b Bullet) Collides(element Element) bool {
	return b.GetHitBox().Intersects(element.GetHitBox())
}

func (b *Bullet) Collision(elements []Element) {
}

func (b Bullet) IsDestroyed() bool {
	return b.destroyed
}
