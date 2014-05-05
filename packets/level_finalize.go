package packets

const LevelFinalizeSize = (ByteSize + // Packet ID
	ShortSize + // X Size
	ShortSize + // Y Size
	ShortSize) // Z Size

type LevelFinalize struct {
	PacketId byte
	X, Y, Z  int16
}

func (p LevelFinalize) Id() byte {
	return p.PacketId
}

func (p LevelFinalize) Size() int {
	return LevelFinalizeSize
}

func (p LevelFinalize) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadLevelFinalize(b []byte) (Packet, error) {
	var p LevelFinalize
	err := ReflectRead(b, &p)
	return &p, err
}

func NewLevelFinalize(x, y, z int16) (p *LevelFinalize, err error) {
	p = &LevelFinalize{
		PacketId: 0x04,
		X:        x,
		Y:        y,
		Z:        z,
	}
	return
}
