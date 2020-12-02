package scenes

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/conf"
	"github.com/lambher/fireblast/models"
)

type Scene struct {
	Conf     *conf.Conf
	Elements []models.Element
	//Player1   *models.Player
	//Player2   *models.Player
	//UIPlayer1 *models.UI
	//UIPlayer2 *models.UI

	Players map[int]*models.Player
	mu      *sync.Mutex

	UIPlayers map[int]*models.UI

	ID int
}

func (s *Scene) Init(conf *conf.Conf) {
	s.Conf = conf
	s.Elements = make([]models.Element, 0)
	s.Players = make(map[int]*models.Player)
	s.UIPlayers = make(map[int]*models.UI)
	s.mu = new(sync.Mutex)
}

//func (s *Scene) AddPlayer1(playerConf models.PlayerConf) {
//	s.Player1 = models.NewPlayer(playerConf)
//	s.Player1.ID = 1
//	s.Elements = append(s.Elements, s.Player1)
//}

func (s Scene) GetPlayer(id int) (*models.Player, bool) {
	s.mu.Lock()
	player, exist := s.Players[id]
	s.mu.Unlock()
	return player, exist
}

func (s *Scene) AddPlayer(playerConf models.PlayerConf, id int) *models.Player {
	player := models.NewPlayer(playerConf)
	player.ID = id
	s.Players[id] = player
	return player
}

func (s *Scene) RemovePlayer(player *models.Player) {
	delete(s.Players, player.ID)
	delete(s.UIPlayers, player.ID)
}

func (s *Scene) AddUIPlayer(id int) {
	x := 1.
	y := s.Conf.MaxY - 11

	switch id {
	case 1:
		x = 1
		y = s.Conf.MaxY - 11
	case 2:
		x = s.Conf.MaxX - 100
		y = s.Conf.MaxY - 11
	case 3:
		x = s.Conf.MaxX - 100
		y = 100
	case 4:
		x = 1
		y = 100
	}

	s.UIPlayers[id] = models.NewUI(s.Players[id], pixel.Vec{
		X: x,
		Y: y,
	})
}

//func (s *Scene) AddUIPlayer1() {
//	s.UIPlayer1 = models.NewUI(s.Player1, pixel.Vec{
//		X: 1,
//		Y: s.Conf.MaxY - 11,
//	})
//}
//
//func (s *Scene) AddUIPlayer2() {
//	s.UIPlayer2 = models.NewUI(s.Player2, pixel.Vec{
//		X: s.Conf.MaxX - 100,
//		Y: s.Conf.MaxY - 11,
//	})
//}
//
//func (s *Scene) AddPlayer2(playerConf models.PlayerConf) {
//	s.Player2 = models.NewPlayer(playerConf)
//	s.Player2.ID = 2
//	s.Elements = append(s.Elements, s.Player2)
//}

func (s *Scene) CatchEvent(win *pixelgl.Window, input chan string) {
	player, exist := s.GetPlayer(s.ID)
	if !exist {
		return
	}

	if win.Pressed(pixelgl.KeyA) {
		input <- fmt.Sprintf("INPUT %d Left", player.ID)
		player.Rotate(0.05)
	}
	if win.Pressed(pixelgl.KeyD) {
		input <- fmt.Sprintf("INPUT %d Right", player.ID)
		player.Rotate(-0.05)
	}
	if win.JustPressed(pixelgl.KeyLeftShift) {
		input <- fmt.Sprintf("INPUT %d GazOn", player.ID)
		player.Gaz(true)
	}
	if win.JustReleased(pixelgl.KeyLeftShift) {
		input <- fmt.Sprintf("INPUT %d GazOff", player.ID)
		player.Gaz(false)
	}
	if win.JustPressed(pixelgl.KeySpace) {
		input <- fmt.Sprintf("INPUT %d Shoot", player.ID)
		bullet := player.Shoot()
		s.AddBullet(bullet)
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

	for _, UIPlayer := range s.UIPlayers {
		UIPlayer.Draw(win)
	}

	for _, player := range s.Players {
		player.Draw(win)
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

	for _, UIPlayer := range s.UIPlayers {
		UIPlayer.Update()
	}
	for _, player := range s.Players {
		player.Update()
		player.Collision(elements)
	}
}
