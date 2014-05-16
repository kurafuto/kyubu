package packets

import "errors"

// SetBlock5 [0x05] is the client only Set Block packet.
// Similar to SetBlock6 [0x06]
type SetBlock5 struct {
	PacketId        byte
	X, Y, Z         int16
	Mode, BlockType byte
}

func (p SetBlock5) Id() byte {
	return 0x05
}

func (p SetBlock5) Size() int {
	return ReflectSize(p)
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

func init() {
	MustRegister(&PacketInfo{
		Id:   0x05,
		Read: ReadSetBlock5,
		Size: ReflectSize(&SetBlock5{}),
		Type: ClientOnly,
		Name: "Set Block [C->S]",
	})
}
