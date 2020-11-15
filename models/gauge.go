package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Gauge struct {
	a        pixel.Vec
	b        pixel.Vec
	c        pixel.Vec
	d        pixel.Vec
	gazColor pixel.RGBA
}

func NewGauge(position pixel.Vec) *Gauge {
	var gauge Gauge

	gauge.a = position.Add(pixel.Vec{
		X: 35,
		Y: -45,
	})
	gauge.b = position.Add(pixel.Vec{
		X: 35,
		Y: -55,
	})
	gauge.gazColor = pixel.RGB(0, 1, 0)

	return &gauge
}

func (g Gauge) Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	imd.Color = g.gazColor
	imd.Push(g.a)
	imd.Color = g.gazColor
	imd.Push(g.b)
	imd.Color = g.gazColor
	imd.Push(g.c)
	imd.Color = g.gazColor
	imd.Push(g.d)
	imd.Polygon(0)

	imd.Draw(win)
}

func (g *Gauge) Update(percentage float64) {
	g.c = g.b.Add(pixel.Vec{
		X: 50 * percentage,
		Y: 0,
	})
	g.d = g.a.Add(pixel.Vec{
		X: 50 * percentage,
		Y: 0,
	})
	g.gazColor.R = 1 - 1*percentage
	g.gazColor.G = 1 * percentage
}
