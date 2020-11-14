package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/models"
)

type Conf struct {
	MaxX     float64
	MaxY     float64
	MinSpeed float64
	MaxSpeed float64
	Nb       int
}

type Scene struct {
	Conf     *Conf
	Players  []*models.Player
	MyPlayer *models.Player
}

func (s *Scene) Init(conf *Conf) {
	s.Conf = conf
	s.Players = make([]*models.Player, 0)

	s.MyPlayer = models.NewPlayer(models.PlayerConf{
		Name: "Lambert",
		Position: pixel.Vec{
			X: conf.MaxX / 2,
			Y: conf.MaxY / 2,
		},
		Color: pixel.RGB(1, 0, 0),
	})
	s.Players = append(s.Players, s.MyPlayer)
}

func (s *Scene) CatchEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
	}

	if win.Pressed(pixelgl.KeyD) {
		s.MyPlayer.Rotate(-0.1)
	}
	if win.Pressed(pixelgl.KeyA) {
		s.MyPlayer.Rotate(0.1)
	}
}

func (s *Scene) Draw(win *pixelgl.Window) {
	for _, player := range s.Players {
		player.Draw(win)
	}
}

func (s *Scene) Update() {
	for _, player := range s.Players {
		player.Update()
	}
}
