package packets

const DisconnectPlayerSize = (ByteSize + // Packet ID
	StringSize) // Disconnect reason

type DisconnectPlayer struct {
	PacketId byte
	Reason   string
}

func (p DisconnectPlayer) Id() byte {
	return p.PacketId
}

func (p DisconnectPlayer) Size() int {
	return DisconnectPlayerSize
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
