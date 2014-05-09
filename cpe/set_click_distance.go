package cpe

import "github.com/sysr-q/kyubu/packets"

const SetClickDistanceSize = (packets.ByteSize + packets.ShortSize)

type SetClickDistance struct {
	PacketId byte
	Distance int16
}

func (p SetClickDistance) Id() byte {
	return p.PacketId
}

func (p SetClickDistance) Size() int {
	return SetClickDistanceSize
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
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x12,
		Read: ReadSetClickDistance,
		Size: SetClickDistanceSize,
		Type: packets.ServerOnly,
		Name: "Set Click Distance (CPE)",
	})
}
