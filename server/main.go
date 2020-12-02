package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/lambher/fireblast/models"

	"github.com/lambher/fireblast/scenes"

	"github.com/tkanos/gonfig"

	"github.com/lambher/fireblast/conf"
)

func main() {
	err := server(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Connexion struct {
	Addr   net.Addr
	Player *models.Player
	Pc     net.PacketConn
}

// server wraps all the UDP echo server functionality.
// ps.: the server is capable of answering to a single
// client at a time.
func server(ctx context.Context) (err error) {
	var c conf.Conf
	connexions := make(map[string]*Connexion)
	err = gonfig.GetConf("./conf.json", &c)
	if err != nil {
		fmt.Println(err)
		return
	}
	var s scenes.Scene
	s.Init(&c)

	// ListenPacket provides us a wrapper around ListenUDP so that
	// we don't need to call `net.ResolveUDPAddr` and then subsequentially
	// perform a `ListenUDP` with the UDP address.
	//
	// The returned value (PacketConn) is pretty much the same as the one
	// from ListenUDP (UDPConn) - the only difference is that `Packet*`
	// methods and interfaces are more broad, also covering `ip`.
	pc, err := net.ListenPacket("udp", c.Address)
	if err != nil {
		return
	}

	// `Close`ing the packet "connection" means cleaning the data structures
	// allocated for holding information about the listening socket.
	defer pc.Close()

	doneChan := make(chan error, 1)
	//buffer := make([]byte, conf.MaxBufferSize)

	playerID := 1
	// Given that waiting for packets to arrive is blocking by nature and we want
	// to be able of canceling such action if desired, we do that in a separate
	// go routine.
	go func() {
		for {
			var connexion *Connexion
			// By reading from the connection into the buffer, we block until there's
			// new content in the socket that we're listening for new packets.
			//
			// Whenever new packets arrive, `buffer` gets filled and we can continue
			// the execution.
			//
			// note.: `buffer` is not being reset between runs.
			//	  It's expected that only `n` reads are read from it whenever
			//	  inspecting its contents.
			buffer := make([]byte, conf.MaxBufferSize)

			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}
			fmt.Println(string(buffer))
			ok := false
			if connexion, ok = connexions[addr.String()]; !ok {
				connexion = &Connexion{
					Addr:   addr,
					Player: nil,
					Pc:     pc,
				}
				connexions[addr.String()] = connexion
				go func() {
					ticker := time.NewTicker(1000 * time.Millisecond / 30)
					for {
						select {
						case <-ticker.C:
							str := ""

							for _, player := range s.Players {
								str += fmt.Sprintf("PLAYER %d DIRECTION_X %f\n", player.ID, player.Shape.Direction.X)
								str += fmt.Sprintf("PLAYER %d DIRECTION_Y %f\n", player.ID, player.Shape.Direction.Y)
								str += fmt.Sprintf("PLAYER %d INERTIE_X %f\n", player.ID, player.Shape.Inertie.X)
								str += fmt.Sprintf("PLAYER %d INERTIE_Y %f\n", player.ID, player.Shape.Inertie.Y)
								str += fmt.Sprintf("PLAYER %d A_X %f\n", player.ID, player.Shape.A.Translation.X)
								str += fmt.Sprintf("PLAYER %d A_Y %f\n", player.ID, player.Shape.A.Translation.Y)
								str += fmt.Sprintf("PLAYER %d B_X %f\n", player.ID, player.Shape.B.Translation.X)
								str += fmt.Sprintf("PLAYER %d B_Y %f\n", player.ID, player.Shape.B.Translation.Y)
								str += fmt.Sprintf("PLAYER %d C_X %f\n", player.ID, player.Shape.C.Translation.X)
								str += fmt.Sprintf("PLAYER %d C_Y %f\n", player.ID, player.Shape.C.Translation.Y)
								str += fmt.Sprintf("PLAYER %d G_X %f\n", player.ID, player.Shape.G.Translation.X)
								str += fmt.Sprintf("PLAYER %d G_Y %f\n", player.ID, player.Shape.G.Translation.Y)
							}
							_, err = pc.WriteTo([]byte(str), addr)
						}
					}
				}()
			}

			data := string(buffer[:n])
			//fmt.Printf("packet-received: bytes=%d from=%s\n",
			//	n, addr.String())
			if strings.HasPrefix(data, "CONNECT") {
				fmt.Println("[" + data + "]")

				var playerConf models.PlayerConf
				playerConf.Conf = &c
				strs := strings.Split(data, " ")
				for i, str := range strs {
					if i == 0 {
						continue
					}
					if i == 1 {
						playerConf.Name = str
					}
					if i == 2 {
						playerConf.Position.X, _ = strconv.ParseFloat(str, 64)
					}
					if i == 3 {
						playerConf.Position.Y, _ = strconv.ParseFloat(str, 64)
					}
				}
				if len(playerConf.Name) > 0 {
					playerConf.Color.R = getRandomFloat(0, 1)
					playerConf.Color.G = getRandomFloat(0, 1)
					playerConf.Color.B = getRandomFloat(0, 1)

					player := s.AddPlayer(playerConf, playerID)
					connexion.Player = player
					connexion.SetID(player)
					for _, player := range s.Players {
						syncAddPlayer(connexions, player)
					}

					playerID++
				}
			}
			if strings.HasPrefix(data, "INPUT") {
				fmt.Println("[" + data + "]")
				strs := strings.Split(data, " ")
				userID := ""
				input := ""
				for i, str := range strs {
					if i == 0 {
						continue
					}
					if i == 1 {
						userID = str
					}
					if i == 2 {
						input = str
					}
				}
				id, err := strconv.Atoi(userID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				player, exist := s.Players[id]
				if !exist {
					continue
				}
				switch input {
				case "Left":
					player.Rotate(0.05)
					syncLeftPlayer(connexions, player)
					fmt.Println("Rotate", player.Shape.Direction)
				case "Right":
					player.Rotate(-0.05)
					syncRightPlayer(connexions, player)
					fmt.Println("Rotate", player.Shape.Direction)
				case "Shoot":
					bullet := player.Shoot()
					s.AddBullet(bullet)
					syncShootPlayer(connexions, player)
					fmt.Println("Shoot")
				case "GazOn":
					player.Gaz(true)
					syncGazOnPlayer(connexions, player)
					fmt.Println("Gaz On")
				case "GazOff":
					player.Gaz(false)
					syncGazOffPlayer(connexions, player)
					fmt.Println("Gaz Off")
				case "Exit":
					s.RemovePlayer(player)
					syncExitPlayer(connexions, player)
					fmt.Println("Exit")
				default:
				}
			}

			// Setting a deadline for the `write` operation allows us to not block
			// for longer than a specific timeout.
			//
			// In the case of a write operation, that'd mean waiting for the send
			// queue to be freed enough so that we are able to proceed.
			deadline := time.Now().Add(conf.TimeOut)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			// Write the packet's contents back to the client.
			n, err = pc.WriteTo(buffer[:n], addr)
			if err != nil {
				doneChan <- err
				return
			}

			//fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	go func() {
		last := time.Now()
		for {
			last = time.Now()
			s.Update()
			dt := time.Since(last)
			time.Sleep(time.Second/60 - dt)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func getRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func syncAddPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.AddPlayer(player)
		}
	}
}

func syncGazOnPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.GazOnPlayer(player)
		}
	}
}

func syncGazOffPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.GazOffPlayer(player)
		}
	}
}

func syncExitPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.ExitPlayer(player)
		}
	}
}

func syncLeftPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.GazLeftPlayer(player)
		}
	}
}

func syncRightPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.GazRightPlayer(player)
		}
	}
}

func syncShootPlayer(connexions map[string]*Connexion, player *models.Player) {
	for _, connexion := range connexions {
		if connexion.Player != nil && connexion.Player != player {
			connexion.ShootPlayer(player)
		}
	}
}

func (c Connexion) SetID(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("ID %d %s %f %f %f", player.ID, player.Name, player.Color.R, player.Color.G, player.Color.B)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) AddPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("NEW_PLAYER %d %s %f %f %f", player.ID, player.Name, player.Color.R, player.Color.G, player.Color.B)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazOnPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER %d GAZ_ON", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazOffPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER %d GAZ_OFF", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) ExitPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER %d EXIT", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazLeftPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER %d LEFT", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazRightPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER %d RIGHT", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) ShootPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER %d SHOOT", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func cleanStr(str string) string {
	newStr := ""
	for _, c := range str {
		if c == 0 {
			break
		}
		newStr += string(c)
	}
	return newStr
}
