package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Triangle struct {
	Player *Player

	Color pixel.RGBA

	A Coordinate
	B Coordinate
	C Coordinate
	G Coordinate

	Direction        pixel.Vec
	Inertie          pixel.Vec
	DirectionInitial pixel.Vec
}

func NewTriangle(player *Player) *Triangle {
	triangle := Triangle{
		Player: player,
		Color:  player.Color,
		A: Coordinate{
			Position: pixel.Vec{
				X: 0,
				Y: 0,
			},
			Translation: player.Position.Position,
		},
		B: Coordinate{
			Position: pixel.Vec{
				X: 10,
				Y: 0,
			},
			Translation: player.Position.Position,
		},
		C: Coordinate{
			Position: pixel.Vec{
				X: 5,
				Y: 10,
			},
			Translation: player.Position.Position,
		},
		G:                Coordinate{},
		Direction:        pixel.Vec{X: 0, Y: 1},
		Inertie:          pixel.Vec{X: 0, Y: 0},
		DirectionInitial: pixel.Vec{X: 0, Y: 1},
	}

	triangle.refreshCenter()
	triangle.A.Position.X -= triangle.G.Position.X
	triangle.A.Position.Y -= triangle.G.Position.Y
	triangle.B.Position.X -= triangle.G.Position.X
	triangle.B.Position.Y -= triangle.G.Position.Y
	triangle.C.Position.X -= triangle.G.Position.X
	triangle.C.Position.Y -= triangle.G.Position.Y
	triangle.refreshCenter()

	return &triangle
}

func (t *Triangle) applyTranslation() {
	t.Translate(t.Inertie)
}

func (t *Triangle) Translate(vec pixel.Vec) {
	t.A.Translation = t.A.Translation.Add(vec)
	t.B.Translation = t.B.Translation.Add(vec)
	t.C.Translation = t.C.Translation.Add(vec)
	t.G.Translation = t.G.Translation.Add(vec)
}

func (t *Triangle) refreshCenter() {
	t.G.Position.X = (t.A.Position.X + t.B.Position.X + t.C.Position.X) / 3
	t.G.Position.Y = (t.A.Position.Y + t.B.Position.Y + t.C.Position.Y) / 3
}

func (t *Triangle) applyRotation() {
	teta := t.Direction.Rotated(-t.DirectionInitial.Angle()).Angle()
	t.rotate(teta)

	t.DirectionInitial = t.Direction
}

func (t *Triangle) checkLimit() {
	if t.G.Position.Add(t.A.Translation).X >= t.Player.Conf.MaxX {
		t.A.Translation.X = 0
		t.B.Translation.X = 0
		t.C.Translation.X = 0
		t.G.Translation.X = 0
	} else if t.G.Position.Add(t.A.Translation).X <= 0 {
		t.A.Translation.X = t.Player.Conf.MaxX
		t.B.Translation.X = t.Player.Conf.MaxX
		t.C.Translation.X = t.Player.Conf.MaxX
		t.G.Translation.X = t.Player.Conf.MaxX
	}

	if t.G.Position.Add(t.A.Translation).Y >= t.Player.Conf.MaxY {
		t.A.Translation.Y = 0
		t.B.Translation.Y = 0
		t.C.Translation.Y = 0
		t.G.Translation.Y = 0
	} else if t.G.Position.Add(t.A.Translation).Y <= 0 {
		t.A.Translation.Y = t.Player.Conf.MaxY
		t.B.Translation.Y = t.Player.Conf.MaxY
		t.C.Translation.Y = t.Player.Conf.MaxY
		t.G.Translation.Y = t.Player.Conf.MaxY
	}
}

func (t *Triangle) rotate(teta float64) {
	t.A.Position = t.A.Position.Rotated(teta)
	t.B.Position = t.B.Position.Rotated(teta)
	t.C.Position = t.C.Position.Rotated(teta)
	t.G.Position = t.G.Position.Rotated(teta)
}

func (t *Triangle) RotateDirection(teta float64) {
	t.Direction = t.Direction.Rotated(teta)
}

func (t Triangle) Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	imd.Color = t.Color
	imd.Push(t.A.Position.Add(t.A.Translation))
	imd.Color = t.Color
	imd.Push(t.B.Position.Add(t.B.Translation))
	imd.Color = t.Color
	imd.Push(t.C.Position.Add(t.C.Translation))
	imd.Polygon(0)

	imd.Draw(win)
}

func (t *Triangle) Update() {
	t.applyTranslation()
	t.refreshCenter()
	t.applyRotation()
	t.checkLimit()
	t.Inertie = t.Inertie.Scaled(0.99)
}

func (t *Triangle) Gaz() {
	if t.Inertie.Eq(pixel.Vec{}) {
		t.Inertie = t.Direction.Scaled(0.1)
	}
	old := t.Inertie
	t.Inertie = t.Inertie.Add(t.Direction.Scaled(0.1))

	if t.Inertie.Len() >= 5 {
		t.Inertie = old
	}
}
