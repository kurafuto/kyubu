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

func (p OrientationUpdate) Id() byte {
	return p.PacketId
}

func (p OrientationUpdate) Size() int {
	return OrientationUpdateSize
}

func (p OrientationUpdate) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadOrientationUpdate(b []byte) (Packet, error) {
	var p OrientationUpdate
	err := ReflectRead(b, &p)
	return &p, err
}

func NewOrientationUpdate(playerId int8, yaw, pitch byte) (p *OrientationUpdate, err error) {
	p = &OrientationUpdate{
		PacketId: 0x0b,
		PlayerId: playerId,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}
