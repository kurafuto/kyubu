package classic

import "github.com/sysr-q/kyubu/packets"

type DespawnPlayer struct {
	PacketId byte
	PlayerId int8
}

func (p DespawnPlayer) Id() byte {
	return 0x0c
}

func (p DespawnPlayer) Size() int {
	return packets.ReflectSize(p)
}

func (p DespawnPlayer) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadDespawnPlayer(b []byte) (packets.Packet, error) {
	var p DespawnPlayer
	err := packets.ReflectRead(b, &p)
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
	packets.Register(&packets.PacketInfo{
		Id:        0x0c,
		Read:      ReadDespawnPlayer,
		Size:      packets.ReflectSize(&DespawnPlayer{}),
		Direction: packets.Anomalous,
		Name:      "DespawnPlayer",
	})
}
