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
	Gauge      *Gauge
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
	ui.Gauge = NewGauge(position)
	return &ui
}

func (u UI) Draw(win *pixelgl.Window) {
	u.basicTxt.Draw(win, pixel.IM)
	u.Gauge.Draw(win)
}

func (u *UI) Update() {
	u.basicTxt.Clear()
	_, err := u.basicTxt.WriteString(fmt.Sprintf("%s\n\nHP: %d\n\nGAZ:", u.Player.Name, u.Player.HP))
	if err != nil {
		fmt.Println(err)
	}

	u.Gauge.Update(u.Player.GazTank)
}

func (u UI) IsDestroyed() bool {
	return false
}
