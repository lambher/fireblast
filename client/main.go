package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/lambher/blop"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lambher/fireblast/conf"
	"github.com/lambher/fireblast/models"
	"github.com/lambher/fireblast/scenes"
	"github.com/tkanos/gonfig"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	var c conf.Conf

	loadSound()

	err := gonfig.GetConf("./conf.json", &c)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := pixelgl.WindowConfig{
		Title:     "FireBlast",
		Bounds:    pixel.R(0, 0, c.MaxX, c.MaxY),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	input := make(chan string)
	var s scenes.Scene

	s.Init(&c)

	go watchInput(input, &s, c.Address)

	input <- "CONNECT Lambert 400 400"

	last := time.Now()
	for !win.Closed() {
		last = time.Now()

		win.Clear(colornames.Black)
		s.CatchEvent(win, input)
		s.Draw(win)
		s.Update()
		win.Update()

		dt := time.Since(last)
		time.Sleep(time.Second/60 - dt)
	}
	input <- fmt.Sprintf("INPUT %d Exit", s.ID)
}

func loadSound() {
	err := blop.LoadSound("hit", "assets/sounds/hit.mp3")
	if err != nil {
		fmt.Println(err)
	}
	err = blop.LoadSound("shot", "assets/sounds/shot.mp3")
	if err != nil {
		fmt.Println(err)
	}
}

func watchInput(input chan string, s *scenes.Scene, address string) {
	r := strings.NewReader("")
	err := client(context.TODO(), address, r, input, s)
	if err != nil {
		panic(err)
	}
}

// client wraps the whole functionality of a UDP client that sends
// a message and waits for a response coming back from the server
// that it initially targetted.
func client(ctx context.Context, address string, reader io.Reader, input chan string, s *scenes.Scene) error {
	// Resolve the UDP address so that we can make use of DialUDP
	// with an actual IP and port instead of a name (in case a
	// hostname is specified).
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	// Although we're not in a connection-oriented transport,
	// the act of `dialing` is analogous to the act of performing
	// a `connect(2)` syscall for a socket of type SOCK_DGRAM:
	// - it forces the underlying socket to only read and write
	//   to and from a specific remote address.
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	// Closes the underlying file descriptor associated with the,
	// socket so that it no longer refers to any file.
	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		for {
			// It is possible that this action blocks, although this
			// should only occur in very resource-intensive situations:
			// - when you've filled up the socket buffer and the OS
			//   can't dequeue the queue fast enough.
			_, err := io.Copy(conn, reader)
			if err != nil {
				doneChan <- err
				return
			}

			buffer := make([]byte, conf.MaxBufferSize)

			// Set a deadline for the ReadOperation so that we don't
			// wait forever for a server that might not respond on
			// a resonable amount of time.
			deadline := time.Now().Add(conf.TimeOut)
			err = conn.SetReadDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			nRead, _, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}
			allDatas := string(buffer[:nRead])
			//fmt.Printf("data size: %d / %d\n", len(allDatas), conf.MaxBufferSize)

			datas := strings.Split(allDatas, "\n")

			for _, data := range datas {
				parseData(s, data)
			}

		}

		//doneChan <- nil
	}()

	go func() {
		for {
			select {
			case value := <-input:
				_, err := conn.Write([]byte(value))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return nil
}

func parseData(s *scenes.Scene, data string) {
	prefix := getPrefix(data)
	switch prefix {
	case "ID":
		fmt.Println(data)
		value, err := getIntAt(data, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		name := getStringAt(data, 2)

		r, err := getFloatAt(data, 3)
		if err != nil {
			fmt.Println(err)
			return
		}
		g, err := getFloatAt(data, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := getFloatAt(data, 5)
		if err != nil {
			fmt.Println(err)
			return
		}

		s.ID = value
		s.AddPlayer(models.PlayerConf{
			Name: name,
			Position: pixel.Vec{
				X: s.Conf.MaxX / 2,
				Y: s.Conf.MaxY / 2,
			},
			Color: pixel.RGB(r, g, b),
			Conf:  s.Conf,
		}, s.ID)
		s.AddUIPlayer(s.ID)
	case "NEW_PLAYER":
		fmt.Println(data)
		addPlayer(s, data)
	case "PLAYER":
		value, err := getIntAt(data, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		player, exist := s.Players[value]
		if !exist {
			return
		}
		refreshPlayer(s, player, data)
	}
}

func refreshPlayer(s *scenes.Scene, player *models.Player, data string) {
	str := getStringAt(data, 2)
	switch str {
	case "EXIT":
		s.RemovePlayer(player)
	case "GAZ_ON":
		fmt.Println(data)
		player.Gaz(true)
	case "GAZ_OFF":
		fmt.Println(data)
		player.Gaz(false)
	case "LEFT":
		fmt.Println(data)
		player.Rotate(.05)
	case "RIGHT":
		fmt.Println(data)
		player.Rotate(-.05)
	case "SHOOT":
		fmt.Println(data)
		bullet := player.Shoot()
		s.AddBullet(bullet)
	case "DIRECTION_X":
		value, err := getFloatAt(data, 3)
		if err == nil {
			player.Shape.Direction.X = value
		}
	case "DIRECTION_Y":
		value, err := getFloatAt(data, 3)
		if err == nil {
			player.Shape.Direction.Y = value
		}
	case "VELOCITY_X":
		value, err := getFloatAt(data, 3)
		if err == nil {
			player.Shape.Velocity.X = value
		}
	case "VELOCITY_Y":
		value, err := getFloatAt(data, 3)
		if err == nil {
			player.Shape.Velocity.Y = value
		}
	case "TRANSLATION_X":
		value, err := getFloatAt(data, 3)
		if err == nil {
			player.Shape.A.Translation.X = value
			player.Shape.B.Translation.X = value
			player.Shape.C.Translation.X = value
			player.Shape.G.Translation.X = value
		}
	case "TRANSLATION_Y":
		value, err := getFloatAt(data, 3)
		if err == nil {
			player.Shape.A.Translation.Y = value
			player.Shape.B.Translation.Y = value
			player.Shape.C.Translation.Y = value
			player.Shape.G.Translation.Y = value
		}
	}

}

func addPlayer(s *scenes.Scene, data string) {
	var playerConf models.PlayerConf
	playerConf.Position.X = s.Conf.MaxX / 2
	playerConf.Position.Y = s.Conf.MaxY / 2

	id, err := getIntAt(data, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := getStringAt(data, 2)

	r, err := getFloatAt(data, 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	g, err := getFloatAt(data, 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := getFloatAt(data, 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	playerConf.Color = pixel.RGB(r, g, b)

	playerConf.Name = name
	playerConf.Conf = s.Conf
	s.AddPlayer(playerConf, id)
	s.AddUIPlayer(id)
}

func getPrefix(str string) string {
	return strings.Split(str, " ")[0]
}

func getStringAt(str string, index int) string {
	strs := strings.Split(str, " ")
	for i, str := range strs {
		if i == index {
			return str
		}
	}
	return ""
}

func getIntAt(str string, index int) (int, error) {
	strs := strings.Split(str, " ")
	for i, str := range strs {
		if i == index {
			return strconv.Atoi(str)
		}
	}
	return 0, errors.New("value not found")
}

func getFloatAt(str string, index int) (float64, error) {
	strs := strings.Split(str, " ")
	for i, str := range strs {
		if i == index {
			value, err := strconv.ParseFloat(str, 64)
			return value, err
		}
	}
	return 0, errors.New("value not found")
}

func getFloat(str string) (float64, error) {
	strs := strings.Split(str, "=")
	if len(strs) == 2 {
		value, err := strconv.ParseFloat(strs[1], 64)
		return value, err
	}
	return 0, errors.New("missing =")
}
