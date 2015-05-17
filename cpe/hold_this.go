package cpe

import "github.com/sysr-q/kyubu/packets"

type HoldThis struct {
	PacketId      byte
	BlockToHold   byte
	PreventChange byte
}

func (p HoldThis) Id() byte {
	return 0x14
}

func (p HoldThis) Size() int {
	return packets.ReflectSize(p)
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
	packets.Register(&packets.PacketInfo{
		Id:        0x14,
		Read:      ReadHoldThis,
		Size:      packets.ReflectSize(&HoldThis{}),
		Direction: packets.Anomalous,
		Name:      "Hold This (CPE)",
	})
}
