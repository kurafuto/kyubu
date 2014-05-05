package client

import (
	"bytes"
	"fmt"
	"github.com/sysr-q/kyubu/auth"
	"github.com/sysr-q/kyubu/chunk"
	"github.com/sysr-q/kyubu/packets"
	"net"
	"time"
)

type Settings struct {
	Server  auth.Server
	Auth    auth.Auth
	Trickle time.Duration
	Debug   bool
}

type Handler func(packets.Packet) error
type State int

const (
	Dead State = iota
	Identification
	ChunkLoad
	Spawning
	Idle
)

type Client struct {
	conn     net.Conn
	settings Settings
	parser   packets.Parser

	Out     chan packets.Packet // C->S
	In      chan packets.Packet // S->C
	pIn     chan packets.Packet // S->C (internal)
	qTicker *time.Ticker

	quit bool

	Map     *chunk.Chunk
	mapRecv bool
	mapData *bytes.Buffer

	ServerName, ServerMOTD string

	Player  *Player
	Players map[int8]*Player

	State    State
	LoggedIn bool
}

func (c *Client) Debugf(s string, v ...interface{}) {
	if !c.settings.Debug {
		return
	}
	fmt.Printf("kyubu: "+s, v...)
}

func (c *Client) Quit() {
	if c.quit {
		return
	}
	c.quit = true
	c.State = Dead
	go func() {
		time.Sleep(1 * time.Second)
		c.conn.Close()
		close(c.Out)
		close(c.In)
	}()
}

func (c *Client) write() {
	defer func() {
		// Sometimes we'll get a panic for receiving on a closed chan,
		// but if we've quit, that's not an issue and can be ignored.
		if err := recover(); err != nil {
			if !c.quit || c.State != Dead {
				panic(err)
			}
		}
	}()

	for !c.quit {
		p := <-c.Out
		if p == nil {
			break
		}

		info, ok := packets.Packets[p.Id()]
		if !ok {
			c.Debugf("unknown packet %#.2x in write queue", p.Id())
			continue
		}

		<-c.qTicker.C // Wait for leaky bucket to drip.

		switch t := p.(type) {
		case *packets.PositionOrientation:
			c.Player.Teleport(POLoc(t))
			// TODO: PositionUpdate, OrientationUpdate, PositionOrientationUpdate
		}

		n, err := c.conn.Write(p.Bytes())
		if err != nil {
			panic(err)
		}
		if n != p.Size() {
			c.Debugf("packet %#.2x, wrote %d bytes, but size = %d", p.Id(), n, p.Size())
		}

		c.Debugf("->[%#.2x] %s: sent!\n", p.Id(), info.Name)
	}
}

func (c *Client) read() {
	for !c.quit {
		packet, err := c.parser.Next()
		if err != nil {
			panic(err) // TODO: panic() on malformed packet?
		}

		// Send to internal queue first. That passes through to c.In.
		c.pIn <- packet
	}
}

func (c *Client) handleLogic() {
	for !c.quit {
		packet := <-c.pIn

		switch p := packet.(type) {
		case *packets.Identification:
			if p.ProtocolVersion != 0x07 {
				panic(fmt.Errorf("kyubu: Received protocol version: %#.2x, expected 0x07", p.ProtocolVersion))
			}
			c.ServerName = p.Name
			c.ServerMOTD = p.KeyMotd
			c.Player.Op = p.UserType == 0x64

		case *packets.LevelInitialize:
			if !c.mapRecv && c.mapData != nil && c.mapData.Len() > 0 {
				panic("kyubu: Received LevelInitialize, but still looking for LevelFinalize")
			}
			c.mapData = bytes.NewBuffer([]byte{})
			c.mapRecv = false
			c.State = ChunkLoad
		case *packets.LevelDataChunk:
			c.mapData.Grow(int(p.ChunkLength))
			c.mapData.Write(p.ChunkData)
		case *packets.LevelFinalize:
			c.mapRecv = true
			ch, err := chunk.Decompress(c.mapData.Bytes(), p.X, p.Y, p.Z)
			if err != nil {
				panic(err)
			}
			c.Map = ch
			c.State = Spawning

		case *packets.SetBlock6:
			c.Map.SetTile(p.X, p.Y, p.Z, p.BlockType)

		case *packets.SpawnPlayer:
			if p.PlayerId == -1 {
				if c.State == Spawning {
					c.State = Idle
					c.LoggedIn = true
				}
				c.Player.Spawn = &Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				}
				c.Player.Location = &Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				}
			} else {
				c.Players[p.PlayerId] = &Player{
					Id:   p.PlayerId,
					Name: p.PlayerName,
					Location: &Loc{
						X:     p.X,
						Y:     p.Y,
						Z:     p.Z,
						Yaw:   p.Yaw,
						Pitch: p.Pitch,
					},
					Spawn: nil,
				}
			}

		case *packets.PositionOrientation:
			if p.PlayerId < 0 {
				c.Player.Teleport(&Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				})
			} else {
				player, ok := c.Players[p.PlayerId]
				if !ok {
					panic("kyubu: received position/orientation for non-existent player")
				}
				player.Teleport(&Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				})
			}

		// TODO: 'queue' updates, rather than simply setting position directly?
		// It's how Notchian client does it, but probably isn't really an issue
		// in this use case.
		case *packets.PositionOrientationUpdate:
			player, ok := c.Players[p.PlayerId]
			if !ok {
				panic("kyubu: received position/orientation update for non-existent player")
			}
			// TODO: These will probably overflow. Account for that.
			player.Location.X += int16(p.X)
			player.Location.Y += int16(p.Y)
			player.Location.Z += int16(p.Z)
			player.Location.Yaw += p.Yaw
			player.Location.Pitch += p.Pitch

		case *packets.PositionUpdate:
			player, ok := c.Players[p.PlayerId]
			if !ok {
				panic("kyubu: received position update for non-existent player")
			}
			// TODO: These will overflow too.
			player.Location.X += int16(p.X)
			player.Location.Y += int16(p.Y)
			player.Location.Z += int16(p.Z)

		case *packets.OrientationUpdate:
			player, ok := c.Players[p.PlayerId]
			if !ok {
				panic("kyubu: received orientation update for non-existent player")
			}
			// TODO: These will overflow too.
			player.Location.Yaw += p.Yaw
			player.Location.Pitch += p.Pitch

		case *packets.DespawnPlayer:
			delete(c.Players, p.PlayerId)

		case *packets.DisconnectPlayer:
			c.Quit()

		case *packets.UpdateUserType:
			c.Player.Op = p.UserType == 0x64
		}

		c.In <- packet // pass it on
	}
}

func New(settings Settings) (c *Client, err error) {
	ident, err := packets.NewIdentification(settings.Auth.Username(), settings.Server.MpPass)
	if err != nil {
		return
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", settings.Server.Address, settings.Server.Port))
	if err != nil {
		return
	}

	c = &Client{
		conn:     conn,
		settings: settings,
		parser:   packets.NewParser(conn),

		Out:     make(chan packets.Packet),
		In:      make(chan packets.Packet),
		pIn:     make(chan packets.Packet),
		qTicker: time.NewTicker(settings.Trickle * time.Millisecond),

		Player:  &Player{Id: -1, Name: settings.Auth.Username()},
		Players: make(map[int8]*Player),

		State: Identification,
	}
	go c.write()
	go c.handleLogic()
	go c.read()

	c.Out <- ident
	return
}
