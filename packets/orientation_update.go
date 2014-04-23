package packets

const OrientationUpdateSize = (ByteSize + // Packet ID
	SByteSize + // Player ID
	ByteSize + // Yaw (Heading)
	ByteSize) // Pitch

type OrientationUpdate struct {
	PacketId   byte
	PlayerId   int8
	Yaw, Pitch byte
}

func (p *OrientationUpdate) Id() byte {
	return p.PacketId
}

func (p *OrientationUpdate) Size() int {
	return PositionOrientationUpdateSize
}

func (p *OrientationUpdate) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	raw.WriteByte(p.Yaw)
	raw.WriteByte(p.Pitch)
	return raw.Buffer.Bytes()
}

func ReadOrientationUpdate(b []byte) (Packet, error) {
	p := OrientationUpdate{}
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

func NewOrientationUpdate(playerId int8, yaw, pitch byte) (p *OrientationUpdate, err error) {
	p = &OrientationUpdate{
		PacketId: 11,
		PlayerId: playerId,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}
