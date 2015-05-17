package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

type ExtInfo struct {
	PacketId       byte
	AppName        string
	ExtensionCount int16
}

func (p ExtInfo) Id() byte {
	return 0x10
}

func (p ExtInfo) Size() int {
	return packets.ReflectSize(p)
}

func (p ExtInfo) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p ExtInfo) String() string {
	return "Negotiation"
}

func ReadExtInfo(b []byte) (packets.Packet, error) {
	var p ExtInfo
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewExtInfo(appName string, extCount int16) (p *ExtInfo, err error) {
	if len(appName) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &ExtInfo{
		PacketId:       0x10,
		AppName:        appName,
		ExtensionCount: extCount,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x10,
		Read:      ReadExtInfo,
		Size:      packets.ReflectSize(&ExtInfo{}),
		Direction: packets.Anomalous,
		Name:      "Ext Info (CPE)",
	})
}
