package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

type EnvSetMapAppearance struct {
	PacketId   byte
	TextureUrl string
	SideBlock  byte
	EdgeBlock  byte
	SideLevel  int16
}

func (p EnvSetMapAppearance) Id() byte {
	return 0x1e
}

func (p EnvSetMapAppearance) Size() int {
	return packets.ReflectSize(p)
}

func (p EnvSetMapAppearance) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p EnvSetMapAppearance) String() string {
	return "EnvMapAppearance"
}

func ReadEnvSetMapAppearance(b []byte) (packets.Packet, error) {
	var p EnvSetMapAppearance
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewEnvSetMapAppearance(textureUrl string, sideBlock, edgeBlock byte, sideLevel int16) (p *EnvSetMapAppearance, err error) {
	if len(textureUrl) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}

	p = &EnvSetMapAppearance{
		PacketId:   0x1e,
		TextureUrl: textureUrl,
		SideBlock:  sideBlock,
		EdgeBlock:  edgeBlock,
		SideLevel:  sideLevel,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x1e,
		Read: ReadEnvSetMapAppearance,
		Size: packets.ReflectSize(&EnvSetMapAppearance{}),
		Type: packets.ServerOnly,
		Name: "Env Set Map Appearance(CPE)",
	})
}
