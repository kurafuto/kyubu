package packets

const LevelFinalizeSize = (ByteSize + // Packet ID
	ShortSize + // X Size
	ShortSize + // Y Size
	ShortSize) // Z Size

type LevelFinalize struct {
	PacketId byte
	X, Y, Z  int16
}

func (p *LevelFinalize) Id() byte {
	return p.PacketId
}

func (p *LevelFinalize) Size() int {
	return LevelFinalizeSize
}

func (p *LevelFinalize) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteShort(p.X)
	raw.WriteShort(p.Y)
	raw.WriteShort(p.Z)
	return raw.Buffer.Bytes()
}

func ReadLevelFinalize(b []byte) (Packet, error) {
	p := LevelFinalize{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if x, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.X = x
	}
	if y, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.Y = y
	}
	if z, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.Z = z
	}
	return &p, nil
}

func NewLevelFinalize(x, y, z int16) (p *LevelFinalize, err error) {
	p = &LevelFinalize{
		PacketId: 4,
		X:        x,
		Y:        y,
		Z:        z,
	}
	return
}
