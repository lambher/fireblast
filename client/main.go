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

	s.AddPlayer1(models.PlayerConf{
		Name: "Lambert",
		Position: pixel.Vec{
			X: c.MaxX / 2,
			Y: c.MaxY / 2,
		},
		Color: pixel.RGB(1, 0, 0),
		Conf:  &c,
	})
	s.AddUIPlayer1()

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
					continue
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
		s.ID = getStringAt(data, 1)
	case "NEW_PLAYER":
		addPlayer(s, data)
	case "PLAYER_1":
		refreshPlayer(s, s.Player1, data)
	case "PLAYER_2":
		refreshPlayer(s, s.Player2, data)
	}

	return

	if strings.HasPrefix(data, "ID") {
		strs := strings.Split(data, " ")
		s.ID = strs[1]
	} else if strings.HasPrefix(data, "NEW_PLAYER") {
		var playerConf models.PlayerConf
		playerConf.Position.X = s.Conf.MaxX / 2
		playerConf.Position.Y = s.Conf.MaxY / 2
		playerConf.Color = pixel.RGB(0, 1, 0)
		strs := strings.Split(data, " ")
		for i, str := range strs {
			if i == 0 {
				continue
			}
			if i == 1 {
				playerConf.Name = str
			}
		}
		playerConf.Conf = s.Conf
		s.AddPlayer2(playerConf)
		s.AddUIPlayer2()
	} else if strings.HasPrefix(data, "PLAYER") {
		nbPlayer := ""
		var player *models.Player

		strs := strings.Split(data, " ")
		for i, str := range strs {
			if i == 0 {
				ss := strings.Split(str, "_")
				if len(ss) == 2 {
					nbPlayer = ss[1]
					switch nbPlayer {
					case "1":
						player = s.Player1
					case "2":
						player = s.Player2
					}
				}
			}
			if i == 1 {
				switch str {
				case "GAZ_ON":
					player.Gaz(true)
				case "GAZ_OFF":
					player.Gaz(false)
				case "LEFT":
					player.Rotate(.05)
				case "RIGHT":
					player.Rotate(-.05)
				case "SHOOT":
					bullet := player.Shoot()
					s.AddBullet(bullet)
				default:
					if strings.HasPrefix(str, "DIRECTION_X") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.Direction.X = value
					} else if strings.HasPrefix(str, "DIRECTION_Y") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.Direction.Y = value
					} else if strings.HasPrefix(str, "INERTIE_X") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.Inertie.X = value
					} else if strings.HasPrefix(str, "INERTIE_Y") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.Inertie.Y = value
					} else if strings.HasPrefix(str, "A_X") {
						value, err := getFloat(str)
						if err != nil {
							fmt.Println(err)
							continue
						}
						player.Shape.A.Translation.X = value
					} else if strings.HasPrefix(str, "A_Y") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.A.Translation.Y = value
					} else if strings.HasPrefix(str, "B_X") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.B.Translation.X = value
					} else if strings.HasPrefix(str, "B_Y") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.B.Translation.Y = value
					} else if strings.HasPrefix(str, "C_X") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.C.Translation.X = value
					} else if strings.HasPrefix(str, "C_Y") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.C.Translation.Y = value
					} else if strings.HasPrefix(str, "G_X") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.G.Translation.X = value
					} else if strings.HasPrefix(str, "G_Y") {
						value, err := getFloat(str)
						if err != nil {
							continue
						}
						player.Shape.G.Translation.Y = value
					}
				}
			}
		}
	}
}

func refreshPlayer(s *scenes.Scene, player *models.Player, data string) {
	str := getStringAt(data, 1)

	switch str {
	case "GAZ_ON":
		player.Gaz(true)
	case "GAZ_OFF":
		player.Gaz(false)
	case "LEFT":
		player.Rotate(.05)
	case "RIGHT":
		player.Rotate(-.05)
	case "SHOOT":
		bullet := player.Shoot()
		s.AddBullet(bullet)
	case "DIRECTION_X":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.Direction.X = value
		}
	case "DIRECTION_Y":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.Direction.Y = value
		}
	case "INERTIE_X":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.Inertie.X = value
		}
	case "INERTIE_Y":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.Inertie.Y = value
		}
	case "A_X":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.A.Translation.X = value
		}
	case "A_Y":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.A.Translation.Y = value
		}
	case "B_X":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.B.Translation.X = value
		}
	case "B_Y":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.B.Translation.Y = value
		}
	case "C_X":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.C.Translation.X = value
		}
	case "C_Y":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.C.Translation.Y = value
		}
	case "G_X":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.G.Translation.X = value
		}
	case "G_Y":
		value, err := getFloatAt(str, 2)
		if err == nil {
			player.Shape.G.Translation.Y = value
		}
	}

}

func addPlayer(s *scenes.Scene, data string) {
	var playerConf models.PlayerConf
	playerConf.Position.X = s.Conf.MaxX / 2
	playerConf.Position.Y = s.Conf.MaxY / 2
	playerConf.Color = pixel.RGB(0, 1, 0)
	strs := strings.Split(data, " ")
	for i, str := range strs {
		if i == 0 {
			continue
		}
		if i == 1 {
			playerConf.Name = str
		}
	}
	playerConf.Conf = s.Conf
	s.AddPlayer2(playerConf)
	s.AddUIPlayer2()
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
