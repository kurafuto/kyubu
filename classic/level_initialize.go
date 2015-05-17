package classic

import "github.com/sysr-q/kyubu/packets"

type LevelInitialize struct {
	PacketId byte
}

func (p LevelInitialize) Id() byte {
	return 0x01
}

func (p LevelInitialize) Size() int {
	return packets.ReflectSize(p)
}

func (p LevelInitialize) Bytes() []byte {
	return []byte{p.PacketId}
}

func ReadLevelInitialize(b []byte) (packets.Packet, error) {
	var p LevelInitialize
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewLevelInitialize() (p *LevelInitialize, err error) {
	p = &LevelInitialize{PacketId: 0x01}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x02,
		Read:      ReadLevelInitialize,
		Size:      packets.ReflectSize(&LevelInitialize{}),
		Direction: packets.Anomalous,
		Name:      "Level Initialize",
	})
}
