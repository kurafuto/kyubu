package packets

const PositionOrientationUpdateSize = (ByteSize + // Packet ID
	SByteSize + // Player ID
	SByteSize + // Change in X
	SByteSize + // Change in Y
	SByteSize + // Change in Z
	ByteSize + // Yaw (Heading)
	ByteSize) // Pitch

type PositionOrientationUpdate struct {
	PacketId   byte
	PlayerId   int8
	X, Y, Z    int8
	Yaw, Pitch byte
}

func (p *PositionOrientationUpdate) Id() byte {
	return p.PacketId
}

func (p *PositionOrientationUpdate) Size() int {
	return PositionOrientationUpdateSize
}

func (p *PositionOrientationUpdate) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	raw.WriteSByte(p.X)
	raw.WriteSByte(p.Y)
	raw.WriteSByte(p.Z)
	raw.WriteByte(p.Yaw)
	raw.WriteByte(p.Pitch)
	return raw.Buffer.Bytes()
}

func ReadPositionOrientationUpdate(b []byte) (Packet, error) {
	p := PositionOrientationUpdate{}
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
	if x, err := raw.ReadSByte(); err != nil {
		return nil, err
	} else {
		p.X = x
	}
	if y, err := raw.ReadSByte(); err != nil {
		return nil, err
	} else {
		p.Y = y
	}
	if z, err := raw.ReadSByte(); err != nil {
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

func NewPositionOrientationUpdate(playerId int8, x, y, z int8, yaw, pitch byte) (p *PositionOrientationUpdate, err error) {
	p = &PositionOrientationUpdate{
		PacketId: 9,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}
