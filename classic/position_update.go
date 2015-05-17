package classic

import "github.com/sysr-q/kyubu/packets"

type PositionUpdate struct {
	PacketId byte
	PlayerId int8
	X, Y, Z  int8
}

func (p PositionUpdate) Id() byte {
	return 0x0a
}

func (p PositionUpdate) Size() int {
	return packets.ReflectSize(p)
}

func (p PositionUpdate) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadPositionUpdate(b []byte) (packets.Packet, error) {
	var p PositionUpdate
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewPositionUpdate(playerId int8, x, y, z int8) (p *PositionUpdate, err error) {
	p = &PositionUpdate{
		PacketId: 0x0a,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x0a,
		Read:      ReadPositionUpdate,
		Size:      packets.ReflectSize(&PositionUpdate{}),
		Direction: packets.Anomalous,
		Name:      "Position Update",
	})
}
