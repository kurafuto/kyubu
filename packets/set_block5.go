package packets

import "errors"

const SetBlock5Size = (ByteSize + // Packet ID
	ShortSize + // X
	ShortSize + // Y
	ShortSize + // Z
	ByteSize + // Mode
	ByteSize) // Block type

// SetBlock5 [0x05] is the client only Set Block packet.
// Similar to SetBlock6 [0x06]
type SetBlock5 struct {
	PacketId        byte
	X, Y, Z         int16
	Mode, BlockType byte
}

func (p SetBlock5) Id() byte {
	return p.PacketId
}

func (p SetBlock5) Size() int {
	return SetBlock5Size
}

func (p SetBlock5) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadSetBlock5(b []byte) (Packet, error) {
	var p SetBlock5
	err := ReflectRead(b, &p)
	return &p, err
}

func NewSetBlock5(x, y, z int16, mode, blockType byte) (p *SetBlock5, err error) {
	if mode != 0 && mode != 1 {
		return nil, errors.New("kyubu/packets: SetBlock mode must be 0x00 or 0x01")
	}
	p = &SetBlock5{
		PacketId:  0x05,
		X:         x,
		Y:         y,
		Z:         z,
		Mode:      mode,
		BlockType: blockType,
	}
	return
}
