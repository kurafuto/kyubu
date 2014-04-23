package packets

const PositionUpdateSize = (ByteSize + // Packet ID
	SByteSize + // Player ID
	SByteSize + // Change in X
	SByteSize + // Change in Y
	SByteSize) // Change in Z

type PositionUpdate struct {
	PacketId byte
	PlayerId int8
	X, Y, Z  int8
}

func (p *PositionUpdate) Id() byte {
	return p.PacketId
}

func (p *PositionUpdate) Size() int {
	return PositionUpdateSize
}

func (p *PositionUpdate) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	raw.WriteSByte(p.X)
	raw.WriteSByte(p.Y)
	raw.WriteSByte(p.Z)
	return raw.Buffer.Bytes()
}

func ReadPositionUpdate(b []byte) (Packet, error) {
	p := PositionUpdate{}
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
	return &p, nil
}

func NewPositionUpdate(playerId int8, x, y, z int8) (p *PositionUpdate, err error) {
	p = &PositionUpdate{
		PacketId: 10,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
	}
	return
}
