package packets

const PositionOrientationSize = (ByteSize + // Packet ID
	SByteSize + // Player ID
	ShortSize + // X
	ShortSize + // Y
	ShortSize + // Z
	ByteSize + // Yaw (Heading)
	ByteSize) // Pitch

type PositionOrientation struct {
	PacketId   byte
	PlayerId   int8
	X, Y, Z    int16
	Yaw, Pitch byte
}

func (p *PositionOrientation) Id() byte {
	return p.PacketId
}

func (p *PositionOrientation) Size() int {
	return PositionOrientationSize
}

func (p *PositionOrientation) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	raw.WriteShort(p.X)
	raw.WriteShort(p.Y)
	raw.WriteShort(p.Z)
	raw.WriteByte(p.Yaw)
	raw.WriteByte(p.Pitch)
	return raw.Buffer.Bytes()
}

func ReadPositionOrientation(b []byte) (Packet, error) {
	p := PositionOrientation{}
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

func NewPositionOrientation(playerId int8, x, y, z int16, yaw, pitch byte) (p *PositionOrientation, err error) {
	p = &PositionOrientation{
		PacketId: 8,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}
