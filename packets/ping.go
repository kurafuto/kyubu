package packets

type Ping struct {
	PacketId byte
}

func (p Ping) Id() byte {
	return 0x01
}

func (p Ping) Size() int {
	return ReflectSize(p)
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

func init() {
	MustRegister(&PacketInfo{
		Id:   0x01,
		Read: ReadPing,
		Size: ReflectSize(&Ping{}),
		Type: ServerOnly,
		Name: "Ping",
	})
}
