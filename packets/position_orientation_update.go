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

func (p PositionOrientationUpdate) Id() byte {
	return p.PacketId
}

func (p PositionOrientationUpdate) Size() int {
	return PositionOrientationUpdateSize
}

func (p PositionOrientationUpdate) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadPositionOrientationUpdate(b []byte) (Packet, error) {
	var p PositionOrientationUpdate
	err := ReflectRead(b, &p)
	return &p, err
}

func NewPositionOrientationUpdate(playerId int8, x, y, z int8, yaw, pitch byte) (p *PositionOrientationUpdate, err error) {
	p = &PositionOrientationUpdate{
		PacketId: 0x09,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}
