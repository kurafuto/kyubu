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

func (p *SpawnPlayer) Id() byte {
	return p.PacketId
}

func (p *SpawnPlayer) Size() int {
	return SpawnPlayerSize
}

func (p *SpawnPlayer) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	raw.WriteString(p.PlayerName)
	raw.WriteShort(p.X)
	raw.WriteShort(p.Y)
	raw.WriteShort(p.Z)
	raw.WriteByte(p.Yaw)
	raw.WriteByte(p.Pitch)
	return raw.Buffer.Bytes()
}

func ReadSpawnPlayer(b []byte) (Packet, error) {
	p := SpawnPlayer{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if playerId, err := raw.ReadSByte(); err != nil {
		return nil, err
	} else {
		p.PlayerId = playerId
	}
	if playerName, err := raw.ReadString(); err != nil {
		return nil, err
	} else {
		p.PlayerName = playerName
	}
	if x, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.X = x
	}
	if y, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.Y = y
	}
	if z, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.Z = z
	}
	if yaw, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.Yaw = yaw
	}
	if pitch, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.Pitch = pitch
	}
	return &p, nil
}

func NewSpawnPlayer(playerId int8, playerName string, x, y, z int16, yaw, pitch byte) (p *SpawnPlayer, err error) {
	if len(playerName) > StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
	}
	p = &SpawnPlayer{
		PacketId:   7,
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
