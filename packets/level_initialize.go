package packets

const LevelInitializeSize = ByteSize // Packet ID

type LevelInitialize struct {
	PacketId byte
}

func (p *LevelInitialize) Id() byte {
	return p.PacketId
}

func (p *LevelInitialize) Size() int {
	return LevelInitializeSize
}

func (p *LevelInitialize) Bytes() []byte {
	return []byte{p.PacketId}
}

func ReadLevelInitialize(b []byte) (Packet, error) {
	p := LevelInitialize{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	return &p, nil
}

func NewLevelInitialize() (p *LevelInitialize, err error) {
	p = &LevelInitialize{PacketId: 1}
	return
}
