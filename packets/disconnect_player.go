package packets

type DisconnectPlayer struct {
	PacketId byte
	Reason   string
}

func (p DisconnectPlayer) Id() byte {
	return 0x0e
}

func (p DisconnectPlayer) Size() int {
	return ReflectSize(p)
}

func (p DisconnectPlayer) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadDisconnectPlayer(b []byte) (Packet, error) {
	var p DisconnectPlayer
	err := ReflectRead(b, &p)
	return &p, err
}

func NewDisconnectPlayer(reason string) (p *DisconnectPlayer, err error) {
	p = &DisconnectPlayer{
		PacketId: 0x0e,
		Reason:   reason,
	}
	return
}

func init() {
	MustRegister(&PacketInfo{
		Id: 0x0e,
		Read: ReadDisconnectPlayer,
		Size: ReflectSize(&DisconnectPlayer{}),
		Type: ServerOnly,
		Name: "Disconnect Player",
	})
}
