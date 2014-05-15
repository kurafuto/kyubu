package cpe

import "github.com/sysr-q/kyubu/packets"

const ExtRemovePlayerNameSize = (packets.ByteSize + packets.ShortSize)

type ExtRemovePlayerName struct {
	PacketId byte
	NameId   int16
}

func (p ExtRemovePlayerName) Id() byte {
	return 0x18
}

func (p ExtRemovePlayerName) Size() int {
	return ExtRemovePlayerNameSize
}

func (p ExtRemovePlayerName) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p ExtRemovePlayerName) String() string {
	return "ExtPlayerList"
}

func ReadExtRemovePlayerName(b []byte) (packets.Packet, error) {
	var p ExtRemovePlayerName
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewExtRemovePlayerName(nameId int16) (p *ExtRemovePlayerName, err error) {
	p = &ExtRemovePlayerName{
		PacketId: 0x18,
		NameId:   nameId,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x18,
		Read: ReadExtRemovePlayerName,
		Size: ExtRemovePlayerNameSize,
		Type: packets.ServerOnly,
		Name: "Ext Remove Player Name (CPE)",
	})
}
