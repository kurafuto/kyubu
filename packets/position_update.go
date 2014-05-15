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

func (p PositionUpdate) Id() byte {
	return 0x0a
}

func (p PositionUpdate) Size() int {
	return PositionUpdateSize
}

func (p PositionUpdate) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadPositionUpdate(b []byte) (Packet, error) {
	var p PositionUpdate
	err := ReflectRead(b, &p)
	return &p, err
}

func NewPositionUpdate(playerId int8, x, y, z int8) (p *PositionUpdate, err error) {
	p = &PositionUpdate{
		PacketId: 0x0a,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
	}
	return
}
