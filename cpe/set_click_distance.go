package cpe

import "github.com/sysr-q/kyubu/packets"

type SetClickDistance struct {
	PacketId byte
	Distance int16
}

func (p SetClickDistance) Id() byte {
	return 0x12
}

func (p SetClickDistance) Size() int {
	return packets.ReflectSize(p)
}

func (p SetClickDistance) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p SetClickDistance) String() string {
	return "ClickDistance"
}

func ReadSetClickDistance(b []byte) (packets.Packet, error) {
	var p SetClickDistance
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewSetClickDistance(distance int16) (p *SetClickDistance, err error) {
	p = &SetClickDistance{
		PacketId: 0x12,
		Distance: distance,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x12,
		Read:      ReadSetClickDistance,
		Size:      packets.ReflectSize(&SetClickDistance{}),
		Direction: packets.Anomalous,
		Name:      "Set Click Distance (CPE)",
	})
}
