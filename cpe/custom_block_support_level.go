package cpe

import "github.com/sysr-q/kyubu/packets"

const CustomBlockSupportLevelSize = (packets.ByteSize + packets.ByteSize)

type CustomBlockSupportLevel struct {
	PacketId     byte
	SupportLevel byte
}

func (p CustomBlockSupportLevel) Id() byte {
	return 0x13
}

func (p CustomBlockSupportLevel) Size() int {
	return CustomBlockSupportLevelSize
}

func (p CustomBlockSupportLevel) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p CustomBlockSupportLevel) String() string {
	return "CustomBlocks"
}

func ReadCustomBlockSupportLevel(b []byte) (packets.Packet, error) {
	var p CustomBlockSupportLevel
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewCustomBlockSupportLevel(supportLevel byte) (p *CustomBlockSupportLevel, err error) {
	p = &CustomBlockSupportLevel{
		PacketId:     0x13,
		SupportLevel: supportLevel,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x13,
		Read: ReadCustomBlockSupportLevel,
		Size: CustomBlockSupportLevelSize,
		Type: packets.Both,
		Name: "Custom Block Support Level (CPE)",
	})
}
