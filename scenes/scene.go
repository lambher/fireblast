package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/conf"
	"github.com/lambher/fireblast/models"
)

type Scene struct {
	Conf     *conf.Conf
	Elements []models.Element
	MyPlayer *models.Player
}

func (s *Scene) Init(conf *conf.Conf) {
	s.Conf = conf
	s.Elements = make([]models.Element, 0)
	s.MyPlayer = models.NewPlayer(models.PlayerConf{
		Conf: conf,
		Name: "Lambert",
		Position: pixel.Vec{
			X: conf.MaxX / 2,
			Y: conf.MaxY / 2,
		},
		Color: pixel.RGB(1, 0, 0),
	})
	s.Elements = append(s.Elements, s.MyPlayer)
}

func (s *Scene) CatchEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
	}

	if win.Pressed(pixelgl.KeyD) {
		s.MyPlayer.Rotate(-0.05)
	}
	if win.Pressed(pixelgl.KeyA) {
		s.MyPlayer.Rotate(0.05)
	}

	if win.JustPressed(pixelgl.KeyLeftShift) {
		s.MyPlayer.Gaz(true)
	}
	if win.JustReleased(pixelgl.KeyLeftShift) {
		s.MyPlayer.Gaz(false)
	}

	if win.JustPressed(pixelgl.KeySpace) {
		bullet := s.MyPlayer.Shoot()
		s.AddBullet(bullet)
	}
}

func (s *Scene) Draw(win *pixelgl.Window) {
	for _, element := range s.Elements {
		element.Draw(win)
	}
}

func (s *Scene) AddBullet(bullet *models.Bullet) {
	s.Elements = append(s.Elements, bullet)
}

func (s *Scene) Update() {
	elements := make([]models.Element, 0)

	for _, element := range s.Elements {
		if !element.IsDestroyed() {
			element.Update()
			elements = append(elements, element)
		}
	}
	s.Elements = elements
}
