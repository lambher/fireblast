package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Element interface {
	Draw(win *pixelgl.Window)
	Update()
	Collision(elements []Element)
	Collides(element Element) bool
	IsDestroyed() bool
	Hit(element Element)
	GetHitBox() pixel.Rect
}
