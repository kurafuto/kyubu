package packets

const PingSize = ByteSize // Packet ID

type Ping struct {
	PacketId byte
}

func (p *Ping) Id() byte {
	return p.PacketId
}

func (p *Ping) Size() int {
	return PingSize
}

func (p *Ping) Bytes() []byte {
	return []byte{p.PacketId}
}

func ReadPing(b []byte) (Packet, error) {
	p := Ping{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	return &p, nil
}

func NewPing() (p *Ping, err error) {
	p = &Ping{PacketId: 1}
	return
}
