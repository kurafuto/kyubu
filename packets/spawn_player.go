package packets

import "fmt"

const SpawnPlayerSize = (ByteSize + // Packet ID
	SByteSize + // Player ID
	StringSize + // Player Name
	ShortSize + // X
	ShortSize + // Y
	ShortSize + // Z
	ByteSize + // Yaw (Heading)
	ByteSize) // Pitch

type SpawnPlayer struct {
	PacketId   byte
	PlayerId   int8
	PlayerName string
	X, Y, Z    int16
	Yaw, Pitch byte
}

func (p SpawnPlayer) Id() byte {
	return p.PacketId
}

func (p SpawnPlayer) Size() int {
	return SpawnPlayerSize
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
