package classic

import "github.com/sysr-q/kyubu/packets"

type OrientationUpdate struct {
	PacketId   byte
	PlayerId   int8
	Yaw, Pitch byte
}

func (p OrientationUpdate) Id() byte {
	return 0x0b
}

func (p OrientationUpdate) Size() int {
	return packets.ReflectSize(p)
}

func (p OrientationUpdate) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadOrientationUpdate(b []byte) (packets.Packet, error) {
	var p OrientationUpdate
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewOrientationUpdate(playerId int8, yaw, pitch byte) (p *OrientationUpdate, err error) {
	p = &OrientationUpdate{
		PacketId: 0x0b,
		PlayerId: playerId,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x0b,
		Read:      ReadOrientationUpdate,
		Size:      packets.ReflectSize(&OrientationUpdate{}),
		Direction: packets.Anomalous,
		Name:      "Orientation Update",
	})
}
