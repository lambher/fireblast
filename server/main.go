package main

import (
	"context"
	"fmt"
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
							str += "PLAYER_1 DIRECTION_X " + fmt.Sprintf("%f\n", s.Player1.Shape.Direction.X)
							str += "PLAYER_1 DIRECTION_Y " + fmt.Sprintf("%f\n", s.Player1.Shape.Direction.Y)
							str += "PLAYER_1 INERTIE_X " + fmt.Sprintf("%f\n", s.Player1.Shape.Inertie.X)
							str += "PLAYER_1 INERTIE_Y " + fmt.Sprintf("%f\n", s.Player1.Shape.Inertie.Y)
							str += "PLAYER_1 A_X " + fmt.Sprintf("%f\n", s.Player1.Shape.A.Translation.X)
							str += "PLAYER_1 A_Y " + fmt.Sprintf("%f\n", s.Player1.Shape.A.Translation.Y)
							str += "PLAYER_1 B_X " + fmt.Sprintf("%f\n", s.Player1.Shape.B.Translation.X)
							str += "PLAYER_1 B_Y " + fmt.Sprintf("%f\n", s.Player1.Shape.B.Translation.Y)
							str += "PLAYER_1 C_X " + fmt.Sprintf("%f\n", s.Player1.Shape.C.Translation.X)
							str += "PLAYER_1 C_Y " + fmt.Sprintf("%f\n", s.Player1.Shape.C.Translation.Y)
							str += "PLAYER_1 G_X " + fmt.Sprintf("%f\n", s.Player1.Shape.G.Translation.X)
							str += "PLAYER_1 G_Y " + fmt.Sprintf("%f\n", s.Player1.Shape.G.Translation.Y)

							if s.Player2 != nil {
								str += "PLAYER_2 DIRECTION_X " + fmt.Sprintf("%f\n", s.Player2.Shape.Direction.X)
								str += "PLAYER_2 DIRECTION_Y " + fmt.Sprintf("%f\n", s.Player2.Shape.Direction.Y)
								str += "PLAYER_2 INERTIE_X " + fmt.Sprintf("%f\n", s.Player2.Shape.Inertie.X)
								str += "PLAYER_2 INERTIE_Y " + fmt.Sprintf("%f\n", s.Player2.Shape.Inertie.Y)
								str += "PLAYER_2 A_X " + fmt.Sprintf("%f\n", s.Player2.Shape.A.Translation.X)
								str += "PLAYER_2 A_Y " + fmt.Sprintf("%f\n", s.Player2.Shape.A.Translation.Y)
								str += "PLAYER_2 B_X " + fmt.Sprintf("%f\n", s.Player2.Shape.B.Translation.X)
								str += "PLAYER_2 B_Y " + fmt.Sprintf("%f\n", s.Player2.Shape.B.Translation.Y)
								str += "PLAYER_2 C_X " + fmt.Sprintf("%f\n", s.Player2.Shape.C.Translation.X)
								str += "PLAYER_2 C_Y " + fmt.Sprintf("%f\n", s.Player2.Shape.C.Translation.Y)
								str += "PLAYER_2 G_X " + fmt.Sprintf("%f\n", s.Player2.Shape.G.Translation.X)
								str += "PLAYER_2 G_Y " + fmt.Sprintf("%f\n", s.Player2.Shape.G.Translation.Y)
							}
							_, err = pc.WriteTo([]byte(str), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 DIRECTION_X="+fmt.Sprintf("%f", s.Player1.Shape.Direction.X)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 DIRECTION_Y="+fmt.Sprintf("%f", s.Player1.Shape.Direction.Y)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 INERTIE_X="+fmt.Sprintf("%f", s.Player1.Shape.Inertie.X)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 INERTIE_Y="+fmt.Sprintf("%f", s.Player1.Shape.Inertie.Y)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 A_X="+fmt.Sprintf("%f", s.Player1.Shape.A.Translation.X)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 A_Y="+fmt.Sprintf("%f", s.Player1.Shape.A.Translation.Y)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 B_X="+fmt.Sprintf("%f", s.Player1.Shape.B.Translation.X)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 B_Y="+fmt.Sprintf("%f", s.Player1.Shape.B.Translation.Y)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 C_X="+fmt.Sprintf("%f", s.Player1.Shape.C.Translation.X)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 C_Y="+fmt.Sprintf("%f", s.Player1.Shape.C.Translation.Y)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 G_X="+fmt.Sprintf("%f", s.Player1.Shape.G.Translation.X)), addr)
							//_, err = pc.WriteTo([]byte("PLAYER_1 G_Y="+fmt.Sprintf("%f", s.Player1.Shape.G.Translation.Y)), addr)
							//
							//if s.Player2 != nil {
							//	_, err = pc.WriteTo([]byte("PLAYER_2 DIRECTION_X="+fmt.Sprintf("%f", s.Player2.Shape.Direction.X)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 DIRECTION_Y="+fmt.Sprintf("%f", s.Player2.Shape.Direction.Y)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 INERTIE_X="+fmt.Sprintf("%f", s.Player2.Shape.Inertie.X)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 INERTIE_Y="+fmt.Sprintf("%f", s.Player2.Shape.Inertie.Y)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 A_X="+fmt.Sprintf("%f", s.Player2.Shape.A.Translation.X)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 A_Y="+fmt.Sprintf("%f", s.Player2.Shape.A.Translation.Y)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 B_X="+fmt.Sprintf("%f", s.Player2.Shape.B.Translation.X)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 B_Y="+fmt.Sprintf("%f", s.Player2.Shape.B.Translation.Y)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 C_X="+fmt.Sprintf("%f", s.Player2.Shape.C.Translation.X)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 C_Y="+fmt.Sprintf("%f", s.Player2.Shape.C.Translation.Y)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 G_X="+fmt.Sprintf("%f", s.Player2.Shape.G.Translation.X)), addr)
							//	_, err = pc.WriteTo([]byte("PLAYER_2 G_Y="+fmt.Sprintf("%f", s.Player2.Shape.G.Translation.Y)), addr)
							//}
						}
					}
				}()
			}

			data := string(buffer[:n])
			//fmt.Printf("packet-received: bytes=%d from=%s\n",
			//	n, addr.String())
			if strings.HasPrefix(data, "CONNECT") {
				fmt.Println("[" + data + "]")

				var player models.PlayerConf
				player.Conf = &c
				strs := strings.Split(data, " ")
				for i, str := range strs {
					if i == 0 {
						continue
					}
					if i == 1 {
						player.Name = str
					}
					if i == 2 {
						player.Position.X, _ = strconv.ParseFloat(str, 64)
					}
					if i == 3 {
						player.Position.Y, _ = strconv.ParseFloat(str, 64)
					}
				}
				if len(player.Name) > 0 {
					if s.Player1 == nil {
						s.AddPlayer1(player)
						connexion.Player = s.Player1
						fmt.Println("ADD Player 1", s.Player1.Name)
						syncAddPlayer(connexions, s.Player1)
						connexion.SetID(s.Player1)
					} else if s.Player2 == nil {
						s.AddPlayer2(player)
						connexion.Player = s.Player2
						fmt.Println("ADD Player 2", s.Player2.Name)
						syncAddPlayer(connexions, s.Player2)
						syncAddPlayer(connexions, s.Player1)
						connexion.SetID(s.Player2)
					}
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
				if userID == "1" && s.Player1 != nil {
					switch input {
					case "Left":
						s.Player1.Rotate(0.05)
						syncLeftPlayer(connexions, s.Player1)
						fmt.Println("Rotate", s.Player1.Shape.Direction)
					case "Right":
						s.Player1.Rotate(-0.05)
						syncRightPlayer(connexions, s.Player1)
						fmt.Println("Rotate", s.Player1.Shape.Direction)
					case "Shoot":
						bullet := s.Player1.Shoot()
						s.AddBullet(bullet)
						syncShootPlayer(connexions, s.Player1)
						fmt.Println("Shoot")
					case "GazOn":
						s.Player1.Gaz(true)
						syncGazOnPlayer(connexions, s.Player1)
						fmt.Println("Gaz On")
					case "GazOff":
						s.Player1.Gaz(false)
						syncGazOffPlayer(connexions, s.Player1)
						fmt.Println("Gaz Off")
					default:
					}
				}
				if userID == "2" && s.Player2 != nil {
					switch input {
					case "Left":
						s.Player2.Rotate(0.05)
						syncLeftPlayer(connexions, s.Player2)
						fmt.Println("Rotate", s.Player2.Shape.Direction)
					case "Right":
						s.Player2.Rotate(-0.05)
						syncRightPlayer(connexions, s.Player2)
						fmt.Println("Rotate", s.Player2.Shape.Direction)
					case "Shoot":
						bullet := s.Player2.Shoot()
						s.AddBullet(bullet)
						syncShootPlayer(connexions, s.Player2)
						fmt.Println("Shoot")
					case "GazOn":
						s.Player2.Gaz(true)
						syncGazOnPlayer(connexions, s.Player2)
						fmt.Println("Gaz On")
					case "GazOff":
						s.Player2.Gaz(false)
						syncGazOffPlayer(connexions, s.Player2)
						fmt.Println("Gaz Off")
					}
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
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("ID %d", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) AddPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte("NEW_PLAYER "+player.Name), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazOnPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER_%d GAZ_ON", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazOffPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER_%d GAZ_OFF", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazLeftPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER_%d LEFT", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) GazRightPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER_%d RIGHT", player.ID)), c.Addr)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Connexion) ShootPlayer(player *models.Player) {
	_, err := c.Pc.WriteTo([]byte(fmt.Sprintf("PLAYER_%d SHOOT", player.ID)), c.Addr)
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
