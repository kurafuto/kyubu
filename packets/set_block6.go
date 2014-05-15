package packets

const SetBlock6Size = (ByteSize + // Packet ID
	ShortSize + // X
	ShortSize + // Y
	ShortSize + // Z
	ByteSize) // Block type

// SetBlock6 [0x06] is the server only Set Block packet.
// Similar to SetBlock5 [0x05]
type SetBlock6 struct {
	PacketId  byte
	X, Y, Z   int16
	BlockType byte
}

func (p SetBlock6) Id() byte {
	return 0x06
}

func (p SetBlock6) Size() int {
	return SetBlock6Size
}

func (p SetBlock6) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadSetBlock6(b []byte) (Packet, error) {
	var p SetBlock6
	err := ReflectRead(b, &p)
	return &p, err
}

func NewSetBlock6(x, y, z int16, blockType byte) (p *SetBlock6, err error) {
	p = &SetBlock6{
		PacketId:  0x06,
		X:         x,
		Y:         y,
		Z:         z,
		BlockType: blockType,
	}
	return
}
