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

func (p *SetBlock6) Id() byte {
	return p.PacketId
}

func (p *SetBlock6) Size() int {
	return SetBlock6Size
}

func (p *SetBlock6) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteShort(p.X)
	raw.WriteShort(p.Y)
	raw.WriteShort(p.Z)
	raw.WriteByte(p.BlockType)
	return raw.Buffer.Bytes()
}

func ReadSetBlock6(b []byte) (Packet, error) {
	p := SetBlock6{}
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
	if blockType, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.BlockType = blockType
	}
	return &p, nil
}

func NewSetBlock6(x, y, z int16, blockType byte) (p *SetBlock6, err error) {
	p = &SetBlock6{
		PacketId:  6,
		X:         x,
		Y:         y,
		Z:         z,
		BlockType: blockType,
	}
	return
}
