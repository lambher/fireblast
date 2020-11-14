package main

import (
	"fmt"

	"github.com/lambher/fireblast/conf"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/scenes"
	"github.com/tkanos/gonfig"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	var conf conf.Conf

	err := gonfig.GetConf("./conf.json", &conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := pixelgl.WindowConfig{
		Title:     "FireBlast",
		Bounds:    pixel.R(0, 0, conf.MaxX, conf.MaxY),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	var s scenes.Scene

	s.Init(&conf)

	for !win.Closed() {
		win.Clear(colornames.Black)
		s.CatchEvent(win)
		s.Draw(win)
		s.Update()
		win.Update()
	}
}
