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

func (p *SetBlock5) Id() byte {
	return p.PacketId
}

func (p *SetBlock5) Size() int {
	return SetBlock5Size
}

func (p *SetBlock5) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteShort(p.X)
	raw.WriteShort(p.Y)
	raw.WriteShort(p.Z)
	raw.WriteByte(p.Mode)
	raw.WriteByte(p.BlockType)
	return raw.Buffer.Bytes()
}

func ReadSetBlock5(b []byte) (Packet, error) {
	p := SetBlock5{}
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
	if mode, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.Mode = mode
	}
	if blockType, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.BlockType = blockType
	}
	return &p, nil
}

func NewSetBlock5(x, y, z int16, mode, blockType byte) (p *SetBlock5, err error) {
	if mode != 0 && mode != 1 {
		return nil, errors.New("kyubu/packets: SetBlock mode must be 0x00 or 0x01")
	}
	p = &SetBlock5{
		PacketId:  5,
		X:         x,
		Y:         y,
		Z:         z,
		Mode:      mode,
		BlockType: blockType,
	}
	return
}
