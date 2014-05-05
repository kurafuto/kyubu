package packets

const PingSize = ByteSize // Packet ID

type Ping struct {
	PacketId byte
}

func (p Ping) Id() byte {
	return p.PacketId
}

func (p Ping) Size() int {
	return PingSize
}

func (p Ping) Bytes() []byte {
	return []byte{p.PacketId}
}

func ReadPing(b []byte) (Packet, error) {
	var p Ping
	err := ReflectRead(b, &p)
	return &p, err
}

func NewPing() (p *Ping, err error) {
	p = &Ping{PacketId: 0x01}
	return
}
