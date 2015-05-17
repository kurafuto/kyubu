package classic

import "github.com/sysr-q/kyubu/packets"

type PositionOrientationUpdate struct {
	PacketId   byte
	PlayerId   int8
	X, Y, Z    int8
	Yaw, Pitch byte
}

func (p PositionOrientationUpdate) Id() byte {
	return 0x09
}

func (p PositionOrientationUpdate) Size() int {
	return packets.ReflectSize(p)
}

func (p PositionOrientationUpdate) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadPositionOrientationUpdate(b []byte) (packets.Packet, error) {
	var p PositionOrientationUpdate
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewPositionOrientationUpdate(playerId int8, x, y, z int8, yaw, pitch byte) (p *PositionOrientationUpdate, err error) {
	p = &PositionOrientationUpdate{
		PacketId: 0x09,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x09,
		Read:      ReadPositionOrientationUpdate,
		Size:      packets.ReflectSize(&PositionOrientationUpdate{}),
		Direction: packets.Anomalous,
		Name:      "Position/Orientation Update",
	})
}
