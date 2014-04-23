package kyubu

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
	"github.com/sysr-q/kyubu/chunk"
	"net"
	"bytes"
)

type Player struct {
	Id int8
	Name string
	Location, Spawn *Loc
	Op bool
}

type Loc struct {
	X, Y, Z int16
	Yaw, Pitch byte
}

type Minecraft struct {
	Host     string
	Port     int
	Conn     net.Conn

	OutQueue chan packets.Packet
	InQueue  chan packets.Packet
	internalQueue chan packets.Packet
	quit     bool

	// Server related
	ServerName, ServerMOTD string
	ServerProtocol byte

	// Map data related
	Map    *chunk.Chunk
	chunkRecvd bool
	chunkData *bytes.Buffer

	// Player related
	ThePlayer *Player
	Players map[byte]*Player

	// Misc
	LoggedIn bool
	waitingForSpawn bool
}

func (m *Minecraft) Quit() {
	m.quit = true
	m.Conn.Close()
	close(m.OutQueue)
	close(m.InQueue)
	close(m.internalQueue)
}

func (m *Minecraft) writeQueue() {
	for {
		if m.quit {
			break
		}
		p := <-m.OutQueue
		_, err := m.Conn.Write(p.Bytes())
		if err != nil {
			panic(err) // debugging
			break
		}
		// I hope you didn't give me a bad packet.
		packetInfo := packets.Packets[p.Id()]
		fmt.Printf("->[%#.2x] %s: sent!\n", p.Id(), packetInfo.Ident)
	}
}

func (m *Minecraft) read() {
	for {
		if m.quit {
			break
		}
		idA := make([]byte, 1)
		n, err := m.Conn.Read(idA)
		if n != 1 {
			panic("kyubu: Server read failed reading fresh packet id")
		} else if err != nil {
			panic(fmt.Errorf("kyubu: Server read failed: %s", err))
		}

		// Packet ID
		id := idA[0]

		packetInfo, ok := packets.Packets[id]
		if !ok {
			panic(fmt.Errorf("kyubu: Unrecognised packet id %#.2x received", id))
		}
		// We've already read ByteSize (1) for the Packet ID
		fixedPacketSize := packetInfo.Size - 1

		packetData := make([]byte, fixedPacketSize)
		if fixedPacketSize != 0 {
			// TODO: This is a hacky mess.
			p := make([]byte, fixedPacketSize)
			n, err = m.Conn.Read(p)
			if n != packetInfo.Size {
				for i := 0; n < fixedPacketSize && i < 5; i++ {
					delta := fixedPacketSize - n // how much we still need
					rest := make([]byte, delta)
					an, err := m.Conn.Read(rest)
					if err != nil {
						break
					}
					n += an
					p = append(p, rest...)
				}
			}
			packetData = append(idA, p...)

			if n != fixedPacketSize {
				panic(fmt.Errorf("kyubu: Packet [%#.2x], wanted to read %d, but got %d", id, fixedPacketSize, n))
			} else if err != nil {
				panic(fmt.Errorf("kyubu: Packet [%#.2x], read failed: %s, wanted %d, got %d", id, err, fixedPacketSize, n))
			}
		} else {
			packetData = append(packetData, id)
		}

		packet, err := packetInfo.Read(packetData)
		if err != nil {
			// TODO: Should malformed packets make us panic() or not?
			panic(err)
		}

		m.InQueue <- packet
		m.internalQueue <- packet
	}
}

func (m *Minecraft) handle() {
	for {
		if m.quit {
			break
		}
		packet := <-m.internalQueue
		switch p := packet.(type) {
			case *packets.Identification:
				if p.ProtocolVersion != 0x07 {
					panic(fmt.Errorf("kyubu: Received protocol version: %#.2x, expected 0x07", p.ProtocolVersion))
				}
				m.ServerName = p.Name
				m.ServerMOTD = p.KeyMotd
				m.ThePlayer.Op = p.UserType == 0x64

			case *packets.LevelInitialize:
				if !m.chunkRecvd && m.chunkData != nil && m.chunkData.Len() > 0 {
					panic("kyubu: Received LevelInitialize, but still looking for LevelFinalize")
				}
				m.chunkData = bytes.NewBuffer([]byte{})
				m.chunkRecvd = false
			case *packets.LevelDataChunk:
				m.chunkData.Grow(int(p.ChunkLength))
				m.chunkData.Write(p.ChunkData)
			case *packets.LevelFinalize:
				m.chunkRecvd = true
				ch, err := chunk.Decompress(m.chunkData.Bytes(), p.X, p.Y, p.Z)
				if err != nil {
					fmt.Printf("%#v\n", p)
					panic(err)
				}
				m.Map = ch
				m.waitingForSpawn = true

			case *packets.SetBlock6:
				m.Map.SetTile(p.X, p.Y, p.Z, p.BlockType)

			case *packets.SpawnPlayer:
				if p.PlayerId == -1 {
					if !m.LoggedIn && m.waitingForSpawn {
						m.waitingForSpawn = false
						m.LoggedIn = true
					}
					m.ThePlayer.Spawn = &Loc{
						X: p.X,
						Y: p.Y,
						Z: p.Z,
						Yaw: p.Yaw,
						Pitch: p.Pitch,
					}
					m.ThePlayer.Location =  &Loc{
						X: p.X,
						Y: p.Y,
						Z: p.Z,
						Yaw: p.Yaw,
						Pitch: p.Pitch,
					}
				} else {
					m.Players[p.PlayerId] = &Player{
						Id: p.PlayerId,
						Name: p.PlayerName,
						Location: &Loc{
							X: p.X,
							Y: p.Y,
							Z: p.Z,
							Yaw: p.Yaw,
							Pitch: p.Pitch,
						},
						Spawn: nil,
					}
				}

			case *packets.PositionOrientation:
				if p.PlayerId < 0 {
					m.ThePlayer.Location.X = p.X
					m.ThePlayer.Location.Y = p.Y
					m.ThePlayer.Location.Z = p.Z
					m.ThePlayer.Location.Yaw = p.Yaw
					m.ThePlayer.Location.Pitch = p.Pitch
				} else {
					player, ok := m.Players[p.PlayerId]
					if !ok {
						panic("kyubu: received position/orientation for non-existent player")
					}
					player.Location.X = p.X
					player.Location.Y = p.Y
					player.Location.Z = p.Z
					player.Location.Yaw = p.Yaw
					player.Location.Pitch = p.Pitch
					if m.Players[p.PlayerId].Location.X != p.X {
						panic("I don't know how Go works.")
					}
				}

			case *packets.DespawnPlayer:
				delete(m.Players, p.PlayerId)
		}
	}
}

func New(host string, port int, username, key string) (m *Minecraft, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return
	}
	m = &Minecraft{
		Host:     host,
		Port:     port,
		Conn:     conn,
		OutQueue: make(chan packets.Packet),
		InQueue:  make(chan packets.Packet),
		internalQueue: make(chan packets.Packet),
		ThePlayer: &Player{Id: -1, Name: username},
		// The rest are zero until otherwise needed.
	}
	go m.writeQueue()
	go m.read()
	go m.handle() // internal

	// Initial identification to server
	p, _ := packets.NewIdentification(username, key)
	m.OutQueue <- p
	return
}
