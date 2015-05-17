package classic

import "github.com/sysr-q/kyubu/packets"

type Ping struct {
	PacketId byte
}

func (p Ping) Id() byte {
	return 0x01
}

func (p Ping) Size() int {
	return packets.ReflectSize(p)
}

func (p Ping) Bytes() []byte {
	return []byte{p.PacketId}
}

func ReadPing(b []byte) (packets.Packet, error) {
	var p Ping
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewPing() (p *Ping, err error) {
	p = &Ping{PacketId: 0x01}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x01,
		Read:      ReadPing,
		Size:      packets.ReflectSize(&Ping{}),
		Direction: packets.Anomalous,
		Name:      "Ping",
	})
}
