package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

const MakeSelectionSize = (packets.ByteSize +
	packets.ByteSize +
	packets.StringSize +
	packets.ShortSize + packets.ShortSize + packets.ShortSize +
	packets.ShortSize + packets.ShortSize + packets.ShortSize +
	packets.ShortSize + packets.ShortSize + packets.ShortSize +
	packets.ShortSize)

type MakeSelection struct {
	PacketId               byte
	SelectionId            byte
	Label                  string
	StartX, StartY, StartZ int16
	EndX, EndY, EndZ       int16
	Red, Green, Blue       int16
	Opacity                int16
}

func (p MakeSelection) Id() byte {
	return 0x1a
}

func (p MakeSelection) Size() int {
	return MakeSelectionSize
}

func (p MakeSelection) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p MakeSelection) String() string {
	return "SelectionCuboid"
}

func ReadMakeSelection(b []byte) (packets.Packet, error) {
	var p MakeSelection
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewMakeSelection(selectionId byte, label string, startX, startY, startZ, endX, endY, endZ, red, green, blue, opacity int16) (p *MakeSelection, err error) {
	if len(label) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &MakeSelection{
		PacketId:    0x1a,
		SelectionId: selectionId,
		Label:       label,
		StartX:      startX,
		StartY:      startY,
		StartZ:      startZ,
		EndX:        endX,
		EndY:        endY,
		EndZ:        endZ,
		Red:         red,
		Green:       green,
		Blue:        blue,
		Opacity:     opacity,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x1a,
		Read: ReadMakeSelection,
		Size: MakeSelectionSize,
		Type: packets.ServerOnly,
		Name: "Make Selection (CPE)",
	})
}
