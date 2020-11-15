package models

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type UI struct {
	Player     *Player
	Position   pixel.Vec
	Color      pixel.RGBA
	basicAtlas *text.Atlas
	basicTxt   *text.Text
	GazGauge   *Gauge
	HPGauge    *Gauge
}

func NewUI(player *Player, position pixel.Vec) *UI {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	ui := UI{
		Player:     player,
		Position:   position,
		basicAtlas: basicAtlas,
		Color:      pixel.RGB(1, 1, 1),
	}
	textPost := position
	ui.basicTxt = text.New(textPost, basicAtlas)
	ui.GazGauge = NewGauge(position.Add(pixel.Vec{
		X: 35,
		Y: -42.5,
	}))
	ui.HPGauge = NewGauge(position.Add(pixel.Vec{
		X: 35,
		Y: -17.5,
	}))
	return &ui
}

func (u UI) Draw(win *pixelgl.Window) {
	u.basicTxt.Draw(win, pixel.IM)
	u.GazGauge.Draw(win)
	u.HPGauge.Draw(win)
}

func (u *UI) Update() {
	u.basicTxt.Clear()
	_, err := u.basicTxt.WriteString(fmt.Sprintf("%s\n\nHP:\n\nGAZ:", u.Player.Name))
	if err != nil {
		fmt.Println(err)
	}

	u.GazGauge.Update(u.Player.GazTank)
	u.HPGauge.Update(float64(u.Player.HP / 100))
}

func (u UI) IsDestroyed() bool {
	return false
}
