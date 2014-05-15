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

func (p PositionOrientation) Id() byte {
	return 0x08
}

func (p PositionOrientation) Size() int {
	return PositionOrientationSize
}

func (p PositionOrientation) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadPositionOrientation(b []byte) (Packet, error) {
	var p PositionOrientation
	err := ReflectRead(b, &p)
	return &p, err
}

func NewPositionOrientation(playerId int8, x, y, z int16, yaw, pitch byte) (p *PositionOrientation, err error) {
	p = &PositionOrientation{
		PacketId: 0x08,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}
