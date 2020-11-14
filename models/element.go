package models

import "github.com/faiface/pixel/pixelgl"

type Element interface {
	Draw(win *pixelgl.Window)
	Update()
	IsDestroyed() bool
}
