package cpe

import "github.com/sysr-q/kyubu/packets"

const RemoveSelectionSize = (packets.ByteSize + packets.ByteSize)

type RemoveSelection struct {
	PacketId    byte
	SelectionId byte
}

func (p RemoveSelection) Id() byte {
	return p.PacketId
}

func (p RemoveSelection) Size() int {
	return RemoveSelectionSize
}

func (p RemoveSelection) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p RemoveSelection) String() string {
	return "SelectionCuboid"
}

func ReadRemoveSelection(b []byte) (packets.Packet, error) {
	var p RemoveSelection
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewRemoveSelection(selectionId byte) (p *RemoveSelection, err error) {
	p = &RemoveSelection{
		PacketId:    0x1b,
		SelectionId: selectionId,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x1b,
		Read: ReadRemoveSelection,
		Size: RemoveSelectionSize,
		Type: packets.ServerOnly,
		Name: "Remove Selection (CPE)",
	})
}
