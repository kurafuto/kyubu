package packets

type DespawnPlayer struct {
	PacketId byte
	PlayerId int8
}

func (p DespawnPlayer) Id() byte {
	return 0x0c
}

func (p DespawnPlayer) Size() int {
	return ReflectSize(p)
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

func init() {
	MustRegister(&PacketInfo{
		Id: 0x0c,
		Read: ReadDespawnPlayer,
		Size: ReflectSize(&DespawnPlayer{}),
		Type: ServerOnly,
		Name: "DespawnPlayer",
	})
}
