package packets

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
	return ReflectSize(p)
}

func (p SpawnPlayer) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadSpawnPlayer(b []byte) (Packet, error) {
	var p SpawnPlayer
	err := ReflectRead(b, &p)
	return &p, err
}

func NewSpawnPlayer(playerId int8, playerName string, x, y, z int16, yaw, pitch byte) (p *SpawnPlayer, err error) {
	if len(playerName) > StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
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
	MustRegister(&PacketInfo{
		Id:   0x07,
		Read: ReadSpawnPlayer,
		Size: ReflectSize(&SpawnPlayer{}),
		Type: ServerOnly,
		Name: "Spawn Player",
	})
}
