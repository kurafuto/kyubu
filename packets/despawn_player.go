package packets

const DespawnPlayerSize = (ByteSize + // Packet ID
	SByteSize) // Player ID

type DespawnPlayer struct {
	PacketId byte
	PlayerId int8
}

func (p *DespawnPlayer) Id() byte {
	return p.PacketId
}

func (p *DespawnPlayer) Size() int {
	return DespawnPlayerSize
}

func (p *DespawnPlayer) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	return raw.Buffer.Bytes()
}

func ReadDespawnPlayer(b []byte) (Packet, error) {
	p := DespawnPlayer{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if playerId, err := raw.ReadSByte(); err != nil {
		return nil, err
	} else {
		p.PlayerId = playerId
	}
	return &p, nil
}

func NewDespawnPlayer(playerId int8) (p *DespawnPlayer, err error) {
	p = &DespawnPlayer{
		PacketId: 12,
		PlayerId: playerId,
	}
	return
}
