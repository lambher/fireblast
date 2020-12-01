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
	ID        string
}

func (s *Scene) Init(conf *conf.Conf) {
	s.Conf = conf
	s.Elements = make([]models.Element, 0)
}

func (s *Scene) AddPlayer1(playerConf models.PlayerConf) {
	s.Player1 = models.NewPlayer(playerConf)
	s.Player1.ID = 1
	s.Elements = append(s.Elements, s.Player1)
}

func (s *Scene) AddUIPlayer1() {
	s.UIPlayer1 = models.NewUI(s.Player1, pixel.Vec{
		X: 1,
		Y: s.Conf.MaxY - 11,
	})
}

func (s *Scene) AddUIPlayer2() {
	s.UIPlayer2 = models.NewUI(s.Player2, pixel.Vec{
		X: s.Conf.MaxX - 100,
		Y: s.Conf.MaxY - 11,
	})
}

func (s *Scene) AddPlayer2(playerConf models.PlayerConf) {
	s.Player2 = models.NewPlayer(playerConf)
	s.Player2.ID = 2
	s.Elements = append(s.Elements, s.Player2)
}

func (s *Scene) CatchEvent(win *pixelgl.Window, input chan string) {
	if s.Player1 != nil && s.ID == "1" {
		if win.Pressed(pixelgl.KeyA) {
			input <- "INPUT 1 Left"
			s.Player1.Rotate(0.05)
		}

		if win.Pressed(pixelgl.KeyD) {
			input <- "INPUT 1 Right"
			s.Player1.Rotate(-0.05)
		}

		if win.JustPressed(pixelgl.KeyLeftShift) {
			input <- "INPUT 1 GazOn"
			s.Player1.Gaz(true)
		}
		if win.JustReleased(pixelgl.KeyLeftShift) {
			input <- "INPUT 1 GazOff"
			s.Player1.Gaz(false)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			input <- "INPUT 1 Shoot"
			bullet := s.Player1.Shoot()
			s.AddBullet(bullet)
		}
	}
	if s.Player2 != nil && s.ID == "2" {
		if win.Pressed(pixelgl.KeyA) {
			input <- "INPUT 2 Left"
			s.Player2.Rotate(0.05)
		}

		if win.Pressed(pixelgl.KeyD) {
			input <- "INPUT 2 Right"
			s.Player2.Rotate(-0.05)
		}

		if win.JustPressed(pixelgl.KeyLeftShift) {
			input <- "INPUT 2 GazOn"
			s.Player2.Gaz(true)
		}
		if win.JustReleased(pixelgl.KeyLeftShift) {
			input <- "INPUT 2 GazOff"
			s.Player2.Gaz(false)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			input <- "INPUT 2 Shoot"
			bullet := s.Player2.Shoot()
			s.AddBullet(bullet)
		}
	}
	//if s.Player2 != nil {
	//	if win.Pressed(pixelgl.KeyLeft) {
	//		s.Player2.Rotate(0.05)
	//		input <- "INPUT 2 Left"
	//	}
	//	if win.Pressed(pixelgl.KeyRight) {
	//		s.Player2.Rotate(-0.05)
	//		input <- "INPUT 2 Right"
	//	}
	//
	//	if win.JustPressed(pixelgl.KeyRightControl) {
	//		s.Player2.Gaz(true)
	//		input <- "INPUT 2 GazOn"
	//	}
	//	if win.JustReleased(pixelgl.KeyRightControl) {
	//		s.Player2.Gaz(false)
	//		input <- "INPUT 2 GazOff"
	//	}
	//
	//	if win.JustPressed(pixelgl.KeyKPEnter) {
	//		bullet := s.Player2.Shoot()
	//		s.AddBullet(bullet)
	//		input <- "INPUT 2 Shoot"
	//	}
	//}
}

func (s *Scene) Draw(win *pixelgl.Window) {
	for _, element := range s.Elements {
		element.Draw(win)
	}

	if s.UIPlayer1 != nil {
		s.UIPlayer1.Draw(win)
	}
	if s.UIPlayer2 != nil {
		s.UIPlayer2.Draw(win)
	}
}

func (s *Scene) AddBullet(bullet *models.Bullet) {
	s.Elements = append(s.Elements, bullet)
}

func (s *Scene) Update() {
	elements := make([]models.Element, 0)

	for _, element := range s.Elements {
		if element != nil && !element.IsDestroyed() {
			element.Update()
			elements = append(elements, element)
		}
	}
	s.Elements = elements

	if s.UIPlayer1 != nil {
		s.UIPlayer1.Update()
	}
	if s.UIPlayer2 != nil {
		s.UIPlayer2.Update()
	}
	if s.Player1 != nil {
		s.Player1.Collision(elements)
	}
	if s.Player2 != nil {
		s.Player2.Collision(elements)
	}
}
