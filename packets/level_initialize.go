package packets

type LevelInitialize struct {
	PacketId byte
}

func (p LevelInitialize) Id() byte {
	return 0x01
}

func (p LevelInitialize) Size() int {
	return ReflectSize(p)
}

func (p LevelInitialize) Bytes() []byte {
	return []byte{p.PacketId}
}

func ReadLevelInitialize(b []byte) (Packet, error) {
	var p LevelInitialize
	err := ReflectRead(b, &p)
	return &p, err
}

func NewLevelInitialize() (p *LevelInitialize, err error) {
	p = &LevelInitialize{PacketId: 0x01}
	return
}

func init() {
	MustRegister(&PacketInfo{
		Id:   0x02,
		Read: ReadLevelInitialize,
		Size: ReflectSize(&LevelInitialize{}),
		Type: ServerOnly,
		Name: "Level Initialize",
	})
}
