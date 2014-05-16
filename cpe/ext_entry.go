package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

type ExtEntry struct {
	PacketId byte
	ExtName  string
	Version  int32
}

func (p ExtEntry) Id() byte {
	return 0x11
}

func (p ExtEntry) Size() int {
	return packets.ReflectSize(p)
}

func (p ExtEntry) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p ExtEntry) String() string {
	return "Negotiation"
}

func ReadExtEntry(b []byte) (packets.Packet, error) {
	var p ExtEntry
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewExtEntry(extName string, version int32) (p *ExtEntry, err error) {
	if len(extName) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &ExtEntry{
		PacketId: 0x11,
		ExtName:  extName,
		Version:  version,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x11,
		Read: ReadExtEntry,
		Size: packets.ReflectSize(&ExtEntry{}),
		Type: packets.Both,
		Name: "Ext Entry (CPE)",
	})
}
