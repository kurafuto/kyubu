package cpe

import "github.com/sysr-q/kyubu/packets"

const HoldThisSize = (packets.ByteSize + packets.ByteSize + packets.ByteSize)

type HoldThis struct {
	PacketId      byte
	BlockToHold   byte
	PreventChange byte
}

func (p HoldThis) Id() byte {
	return p.PacketId
}

func (p HoldThis) Size() int {
	return HoldThisSize
}

func (p HoldThis) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p HoldThis) String() string {
	return "HeldBlock"
}

func ReadHoldThis(b []byte) (packets.Packet, error) {
	var p HoldThis
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewHoldThis(blockToHold, preventChange byte) (p *HoldThis, err error) {
	p = &HoldThis{
		PacketId:      0x14,
		BlockToHold:   blockToHold,
		PreventChange: preventChange,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x14,
		Read: ReadHoldThis,
		Size: HoldThisSize,
		Type: packets.ServerOnly,
		Name: "Hold This (CPE)",
	})
}
