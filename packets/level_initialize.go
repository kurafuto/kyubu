package packets

const LevelInitializeSize = ByteSize // Packet ID

type LevelInitialize struct {
	PacketId byte
}

func (p LevelInitialize) Id() byte {
	return 0x01
}

func (p LevelInitialize) Size() int {
	return LevelInitializeSize
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
