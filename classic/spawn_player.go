package classic

import "github.com/sysr-q/kyubu/packets"

import "fmt"

type SpawnPlayer struct {
	PacketId   byte
	PlayerId   int8
	PlayerName string
	X, Y, Z    int16
	Yaw, Pitch byte
}

func (p SpawnPlayer) Id() byte {
	return 0x07
}

func (p SpawnPlayer) Size() int {
	return packets.ReflectSize(p)
}

func (p SpawnPlayer) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadSpawnPlayer(b []byte) (packets.Packet, error) {
	var p SpawnPlayer
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewSpawnPlayer(playerId int8, playerName string, x, y, z int16, yaw, pitch byte) (p *SpawnPlayer, err error) {
	if len(playerName) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &SpawnPlayer{
		PacketId:   0x07,
		PlayerId:   playerId,
		PlayerName: playerName,
		X:          x,
		Y:          y,
		Z:          z,
		Yaw:        yaw,
		Pitch:      pitch,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x07,
		Read:      ReadSpawnPlayer,
		Size:      packets.ReflectSize(&SpawnPlayer{}),
		Direction: packets.Anomalous,
		Name:      "Spawn Player",
	})
}
