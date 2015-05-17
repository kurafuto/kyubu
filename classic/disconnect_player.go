package classic

import "github.com/sysr-q/kyubu/packets"

type DisconnectPlayer struct {
	PacketId byte
	Reason   string
}

func (p DisconnectPlayer) Id() byte {
	return 0x0e
}

func (p DisconnectPlayer) Size() int {
	return packets.ReflectSize(p)
}

func (p DisconnectPlayer) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadDisconnectPlayer(b []byte) (packets.Packet, error) {
	var p DisconnectPlayer
	err := packets.ReflectRead(b, &p)
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
	packets.Register(&packets.PacketInfo{
		Id:        0x0e,
		Read:      ReadDisconnectPlayer,
		Size:      packets.ReflectSize(&DisconnectPlayer{}),
		Direction: packets.Anomalous,
		Name:      "Disconnect Player",
	})
}
