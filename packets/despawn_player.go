package packets

const DespawnPlayerSize = (ByteSize + // Packet ID
	SByteSize) // Player ID

type DespawnPlayer struct {
	PacketId byte
	PlayerId int8
}

func (p DespawnPlayer) Id() byte {
	return p.PacketId
}

func (p DespawnPlayer) Size() int {
	return DespawnPlayerSize
}

func (p DespawnPlayer) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadDespawnPlayer(b []byte) (Packet, error) {
	var p DespawnPlayer
	err := ReflectRead(b, &p)
	return &p, err
}

func NewDespawnPlayer(playerId int8) (p *DespawnPlayer, err error) {
	p = &DespawnPlayer{
		PacketId: 0x0c,
		PlayerId: playerId,
	}
	return
}
