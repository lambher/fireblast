package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/blop"
	"github.com/lambher/fireblast/conf"
)

type Player struct {
	Conf *conf.Conf

	ID         int
	Name       string
	HP         int
	GazTank    float64
	Position   *Coordinate
	Color      pixel.RGBA
	Invincible bool
	Shape      *Triangle
	gazActive  bool
}

type PlayerConf struct {
	Name     string
	Position pixel.Vec
	Color    pixel.RGBA
	Conf     *conf.Conf
}

func NewPlayer(conf PlayerConf) *Player {
	player := Player{
		Conf:       conf.Conf,
		Name:       conf.Name,
		HP:         100,
		GazTank:    1,
		Invincible: false,
		Position:   NewCoordinate(conf.Position),
		Color:      conf.Color,
	}

	player.Shape = NewTriangle(&player)

	return &player
}

func (p Player) Draw(win *pixelgl.Window) {
	if p.HP <= 0 {
		return
	}

	p.Shape.Draw(win)

	//hitBox := p.GetHitBox()
	//
	//imd := imdraw.New(nil)
	//imd.Color = pixel.RGB(1, 0, 0)
	//imd.Push(p.Shape.G.Translation.Add(p.Shape.G.Position.Add(pixel.V(hitBox.Min.X, 0))))
	//imd.Color = pixel.RGB(1, 0, 0)
	//imd.Push(p.Shape.G.Translation.Add(p.Shape.G.Position.Add(pixel.V(0, hitBox.Min.Y))))
	//imd.Color = pixel.RGB(0, 1, 0)
	//imd.Push(p.Shape.G.Translation.Add(p.Shape.G.Position.Add(pixel.V(hitBox.Max.X, 0))))
	//imd.Color = pixel.RGB(0, 1, 0)
	//imd.Push(p.Shape.G.Translation.Add(p.Shape.G.Position.Add(pixel.V(0, hitBox.Max.Y))))
	////imd.Color = pixel.RGB(0, 0, 1)
	////imd.Push(pixel.V(500, 700))
	//imd.Polygon(0)
	//imd.Draw(win)
}

func (p *Player) Update() {
	if p.gazActive && p.GazTank > 0 {
		p.Shape.Gaz()
		p.GazTank -= 0.01
	} else if !p.gazActive && p.GazTank < 1 {
		p.GazTank += 0.005
		if p.GazTank > 1 {
			p.GazTank = 1
		}
	}
	p.Shape.Update()
}

func (p *Player) Collision(elements []Element) {
	for _, element := range elements {
		bullet := element.(*Bullet)
		if bullet.Player.ID == p.ID {
			continue
		}
		bulletPos := bullet.Position.Position.Add(bullet.Position.Translation)
		if !p.Invincible && bulletPos.X >= p.Shape.G.GetPos().X-10 && bulletPos.X <= p.Shape.G.GetPos().X+10 && bulletPos.Y >= p.Shape.G.GetPos().Y-10 && bulletPos.Y <= p.Shape.G.GetPos().Y+10 {
			p.HP -= 10
			bullet.destroyed = true
			blop.Play("hit")
			continue
		}
	}
}

func (p *Player) Hit(element Element) {

}

func (p Player) Collides(element Element) bool {
	return p.GetHitBox().Intersects(element.GetHitBox())
}

func (p *Player) Shoot() *Bullet {
	blop.Play("shot")
	bullet := NewBullet(p.Conf, p)
	bullet.Position = NewCoordinate(p.Shape.A.Translation)
	bullet.Velocity = p.Shape.Direction.Scaled(5).Add(p.Shape.Velocity).Scaled(1.5)
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

func (p Player) GetHitBox() pixel.Rect {
	return pixel.Rect{
		Min: p.Shape.G.Translation.Add(p.Shape.G.Position.Add(pixel.Vec{-5, -5})),
		Max: p.Shape.G.Translation.Add(p.Shape.G.Position.Add(pixel.Vec{5, 5})),
	}
}
