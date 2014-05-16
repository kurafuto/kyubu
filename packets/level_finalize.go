package packets

type LevelFinalize struct {
	PacketId byte
	X, Y, Z  int16
}

func (p LevelFinalize) Id() byte {
	return 0x04
}

func (p LevelFinalize) Size() int {
	return ReflectSize(p)
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

func init() {
	MustRegister(&PacketInfo{
		Id:   0x04,
		Read: ReadLevelFinalize,
		Size: ReflectSize(&LevelFinalize{}),
		Type: ServerOnly,
		Name: "Level Finalize",
	})
}
