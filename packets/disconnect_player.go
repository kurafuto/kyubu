package packets

const DisconnectPlayerSize = (ByteSize + // Packet ID
	StringSize) // Disconnect reason

type DisconnectPlayer struct {
	PacketId byte
	Reason   string
}

func (p *DisconnectPlayer) Id() byte {
	return p.PacketId
}

func (p *DisconnectPlayer) Size() int {
	return DisconnectPlayerSize
}

func (p *DisconnectPlayer) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteString(p.Reason)
	return raw.Buffer.Bytes()
}

func ReadDisconnectPlayer(b []byte) (Packet, error) {
	p := DisconnectPlayer{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if reason, err := raw.ReadString(); err != nil {
		return nil, err
	} else {
		p.Reason = reason
	}
	return &p, nil
}

func NewDisconnectPlayer(reason string) (p *DisconnectPlayer, err error) {
	p = &DisconnectPlayer{
		PacketId: 14,
		Reason:   reason,
	}
	return
}
