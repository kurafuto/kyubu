package cpe

import "github.com/sysr-q/kyubu/packets"

const EnvSetColorSize = (packets.ByteSize + packets.ByteSize + packets.ShortSize + packets.ShortSize + packets.ShortSize)

type EnvSetColor struct {
	PacketId         byte
	Variable         byte
	Red, Green, Blue int16
}

func (p EnvSetColor) Id() byte {
	return 0x19
}

func (p EnvSetColor) Size() int {
	return EnvSetColorSize
}

func (p EnvSetColor) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p EnvSetColor) String() string {
	return "EnvColors"
}

func ReadEnvSetColor(b []byte) (packets.Packet, error) {
	var p EnvSetColor
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewEnvSetColor(variable byte, red, green, blue int16) (p *EnvSetColor, err error) {
	p = &EnvSetColor{
		PacketId: 0x19,
		Variable: variable,
		Red:      red,
		Green:    green,
		Blue:     blue,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x19,
		Read: ReadEnvSetColor,
		Size: EnvSetColorSize,
		Type: packets.ServerOnly,
		Name: "Env Set Color (CPE)",
	})
}
