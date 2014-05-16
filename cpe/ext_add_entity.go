package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

type ExtAddEntity struct {
	PacketId   byte
	EntityId   byte
	InGameName string
	SkinName   string
}

func (p ExtAddEntity) Id() byte {
	return 0x17
}

func (p ExtAddEntity) Size() int {
	return packets.ReflectSize(p)
}

func (p ExtAddEntity) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p ExtAddEntity) String() string {
	return "ExtPlayerList"
}

func ReadExtAddEntity(b []byte) (packets.Packet, error) {
	var p ExtAddEntity
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewExtAddEntity(entityId byte, inGameName, skinName string) (p *ExtAddEntity, err error) {
	if len(inGameName) > packets.StringSize || len(skinName) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &ExtAddEntity{
		PacketId:   0x17,
		EntityId:   entityId,
		InGameName: inGameName,
		SkinName:   skinName,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x17,
		Read: ReadExtAddEntity,
		Size: packets.ReflectSize(&ExtAddEntity{}),
		Type: packets.ServerOnly,
		Name: "Ext Add Entity (CPE)",
	})
}
