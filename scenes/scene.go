package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/conf"
	"github.com/lambher/fireblast/models"
)

type Scene struct {
	Conf      *conf.Conf
	Elements  []models.Element
	Player1   *models.Player
	Player2   *models.Player
	UIPlayer1 *models.UI
	UIPlayer2 *models.UI
}

func (s *Scene) Init(conf *conf.Conf) {
	s.Conf = conf
	s.Elements = make([]models.Element, 0)
	s.Player1 = models.NewPlayer(models.PlayerConf{
		Conf: conf,
		Name: "Lambert",
		Position: pixel.Vec{
			X: conf.MaxX / 2,
			Y: conf.MaxY / 2,
		},
		Color: pixel.RGB(1, 0, 0),
	})
	s.Player2 = models.NewPlayer(models.PlayerConf{
		Conf: conf,
		Name: "Milande",
		Position: pixel.Vec{
			X: conf.MaxX / 2,
			Y: conf.MaxY / 2,
		},
		Color: pixel.RGB(0, 1, 0),
	})
	s.Elements = append(s.Elements, s.Player1)
	s.Elements = append(s.Elements, s.Player2)

	s.UIPlayer1 = models.NewUI(s.Player1, pixel.Vec{
		X: 1,
		Y: conf.MaxY - 11,
	})

	s.UIPlayer2 = models.NewUI(s.Player2, pixel.Vec{
		X: conf.MaxX - 100,
		Y: conf.MaxY - 11,
	})
}

func (s *Scene) CatchEvent(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyA) {
		s.Player1.Rotate(0.05)
	}

	if win.Pressed(pixelgl.KeyD) {
		s.Player1.Rotate(-0.05)
	}

	if win.JustPressed(pixelgl.KeyLeftShift) {
		s.Player1.Gaz(true)
	}
	if win.JustReleased(pixelgl.KeyLeftShift) {
		s.Player1.Gaz(false)
	}

	if win.JustPressed(pixelgl.KeySpace) {
		bullet := s.Player1.Shoot()
		s.AddBullet(bullet)
	}

	if win.Pressed(pixelgl.KeyLeft) {
		s.Player2.Rotate(0.05)
	}
	if win.Pressed(pixelgl.KeyRight) {
		s.Player2.Rotate(-0.05)
	}

	if win.JustPressed(pixelgl.KeyRightControl) {
		s.Player2.Gaz(true)
	}
	if win.JustReleased(pixelgl.KeyRightControl) {
		s.Player2.Gaz(false)
	}

	if win.JustPressed(pixelgl.KeyKPEnter) {
		bullet := s.Player2.Shoot()
		s.AddBullet(bullet)
	}
}

func (s *Scene) Draw(win *pixelgl.Window) {
	for _, element := range s.Elements {
		element.Draw(win)
	}

	s.UIPlayer1.Draw(win)
	s.UIPlayer2.Draw(win)
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

	s.UIPlayer1.Update()
	s.UIPlayer2.Update()
}
