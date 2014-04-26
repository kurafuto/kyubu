package kyubu

import (
	"bytes"
	"fmt"
	"github.com/sysr-q/kyubu/auth"
	"github.com/sysr-q/kyubu/chunk"
	"github.com/sysr-q/kyubu/packets"
	"net"
	"time"
	"regexp"
	"errors"
	"strconv"
)

type Settings struct {
	Server  auth.Server
	Auth    auth.Auth
	Trickle time.Duration
	Debug   bool
}

var directRegexp = regexp.MustCompile(`^mc://(?P<host>[a-z0-9.-]+):(?P<port>\d+)/(?P<user>[a-zA-Z0-9_]{1,16})/(?P<hash>[a-fA-F0-9]{32})$`)
func Direct(url string) (s Settings, err error) {
	direct := directRegexp.FindStringSubmatch(url)
	if direct == nil {
		err = errors.New("kyubu: Failed to parse direct connect URL")
		return
	}
	port, err := strconv.Atoi(direct[2])
	if err != nil {
		return
	}
	s = Settings{
		Server: auth.Server{Address: direct[1], Port: port, MpPass: direct[4]},
		Auth: auth.NewDirect(direct[3], ""),
		Trickle: 25,
		Debug: false,
	}
	return
}

type Kyubu struct {
	Conn     net.Conn
	settings Settings
	parser   packets.Parser

	Out       chan packets.Packet
	outTicker *time.Ticker
	In        chan packets.Packet
	inAll     chan packets.Packet
	quit      bool

	ServerName, ServerMOTD string
	ServerProtocol         byte

	Map        *chunk.Chunk
	chunkRecvd bool
	chunkData  *bytes.Buffer

	ThePlayer *Player
	Players   map[int8]*Player

	LoggedIn        bool
	waitingForSpawn bool
}

func (k *Kyubu) Debugf(s string, v ...interface{}) {
	if !k.settings.Debug {
		return
	}
	fmt.Printf("kyubu: "+s, v...)
}

func (k *Kyubu) Quit() {
	k.quit = true
	k.LoggedIn = false
	go func() {
		time.Sleep(1000 * time.Millisecond)
		k.Conn.Close()
		close(k.Out)
		close(k.In)
		close(k.inAll)
	}()
}

func (k *Kyubu) writeQueue() {
	defer func() {
		if err := recover(); err != nil {
			if !k.quit || k.LoggedIn {
				panic(err)
			}
		}
	}()

	for !k.quit {
		p := <-k.Out
		<-k.outTicker.C
		_, err := k.Conn.Write(p.Bytes())
		if err != nil {
			panic(err)
		}

		switch t := p.(type) {
		case *packets.PositionOrientation:
			k.Debugf("  [TELE] %#v\n", t)
			k.ThePlayer.Teleport(&Loc{
				X:     t.X,
				Y:     t.Y,
				Z:     t.Z,
				Yaw:   t.Yaw,
				Pitch: t.Pitch,
			})
		}
		packetInfo := packets.Packets[p.Id()]
		k.Debugf("->[%#.2x] %s: sent!\n", p.Id(), packetInfo.Ident)
	}
}

func (k *Kyubu) read() {
	for !k.quit {
		packet, err := k.parser.Next()
		if err != nil {
			// TODO: should malformed packets make us panic() out?
			panic(err)
		}

		k.In <- packet
		k.inAll <- packet
	}
}

func (k *Kyubu) handle() {
	for !k.quit {
		packet := <-k.inAll

		switch p := packet.(type) {
		case *packets.Identification:
			if p.ProtocolVersion != 0x07 {
				panic(fmt.Errorf("kyubu: Received protocol version: %#.2x, expected 0x07", p.ProtocolVersion))
			}
			k.ServerName = p.Name
			k.ServerMOTD = p.KeyMotd
			k.ThePlayer.Op = p.UserType == 0x64

		case *packets.LevelInitialize:
			if !k.chunkRecvd && k.chunkData != nil && k.chunkData.Len() > 0 {
				panic("kyubu: Received LevelInitialize, but still looking for LevelFinalize")
			}
			k.chunkData = bytes.NewBuffer([]byte{})
			k.chunkRecvd = false
		case *packets.LevelDataChunk:
			k.chunkData.Grow(int(p.ChunkLength))
			k.chunkData.Write(p.ChunkData)
		case *packets.LevelFinalize:
			k.chunkRecvd = true
			ch, err := chunk.Decompress(k.chunkData.Bytes(), p.X, p.Y, p.Z)
			if err != nil {
				panic(err)
			}
			k.Map = ch
			k.waitingForSpawn = true

		case *packets.SetBlock6:
			k.Map.SetTile(p.X, p.Y, p.Z, p.BlockType)

		case *packets.SpawnPlayer:
			if p.PlayerId == -1 {
				if !k.LoggedIn && k.waitingForSpawn {
					k.waitingForSpawn = false
					k.LoggedIn = true
				}
				k.ThePlayer.Spawn = &Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				}
				k.ThePlayer.Location = &Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				}
			} else {
				k.Players[p.PlayerId] = &Player{
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
				k.ThePlayer.Teleport(&Loc{
					X:     p.X,
					Y:     p.Y,
					Z:     p.Z,
					Yaw:   p.Yaw,
					Pitch: p.Pitch,
				})
			} else {
				player, ok := k.Players[p.PlayerId]
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
			player, ok := k.Players[p.PlayerId]
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
			player, ok := k.Players[p.PlayerId]
			if !ok {
				panic("kyubu: received position update for non-existent player")
			}
			// TODO: These will overflow too.
			player.Location.X += int16(p.X)
			player.Location.Y += int16(p.Y)
			player.Location.Z += int16(p.Z)

		case *packets.OrientationUpdate:
			player, ok := k.Players[p.PlayerId]
			if !ok {
				panic("kyubu: received orientation update for non-existent player")
			}
			// TODO: These will overflow too.
			player.Location.Yaw += p.Yaw
			player.Location.Pitch += p.Pitch

		case *packets.DespawnPlayer:
			delete(k.Players, p.PlayerId)

		case *packets.DisconnectPlayer:
			k.Quit()

		case *packets.UpdateUserType:
			k.ThePlayer.Op = p.UserType == 0x64
		}
	}
}

// New initializes an instance of *kyubu.Minecraft, starting the session.
func New(settings Settings) (k *Kyubu, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", settings.Server.Address, settings.Server.Port))
	if err != nil {
		return
	}
	k = &Kyubu{
		Conn:     conn,
		settings: settings,
		parser:   packets.NewParser(conn),
		Out:      make(chan packets.Packet),
		// Trickle like a steadily leaking bucket.
		outTicker: time.NewTicker(settings.Trickle * time.Millisecond),
		In:        make(chan packets.Packet),
		inAll:     make(chan packets.Packet),
		ThePlayer: &Player{Id: -1, Name: settings.Auth.Username()},
		Players:   make(map[int8]*Player),
	}
	go k.writeQueue()
	go k.read()
	go k.handle()

	// Initial identification to server
	p, _ := packets.NewIdentification(settings.Auth.Username(), settings.Server.MpPass)
	k.Out <- p
	return
}
