package packets

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
	return ReflectSize(p)
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

func init() {
	MustRegister(&PacketInfo{
		Id:   0x06,
		Read: ReadSetBlock6,
		Size: ReflectSize(&SetBlock6{}),
		Type: ServerOnly,
		Name: "Set Block [S->C]",
	})
}
