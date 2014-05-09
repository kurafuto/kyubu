package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

const ExtAddPlayerNameSize = (packets.ByteSize + packets.ShortSize + packets.StringSize + packets.StringSize + packets.StringSize + packets.ByteSize)

type ExtAddPlayerName struct {
	PacketId   byte
	NameId     int16
	PlayerName string
	ListName   string
	GroupName  string
	GroupRank  byte
}

func (p ExtAddPlayerName) Id() byte {
	return p.PacketId
}

func (p ExtAddPlayerName) Size() int {
	return ExtAddPlayerNameSize
}

func (p ExtAddPlayerName) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p ExtAddPlayerName) String() string {
	return "ExtPlayerList"
}

func ReadExtAddPlayerName(b []byte) (packets.Packet, error) {
	var p ExtAddPlayerName
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewExtAddPlayerName(nameId int16, playerName, listName, groupName string, groupRank byte) (p *ExtAddPlayerName, err error) {
	if len(playerName) > packets.StringSize || len(listName) > packets.StringSize || len(groupName) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &ExtAddPlayerName{
		PacketId:   0x16,
		NameId:     nameId,
		PlayerName: playerName,
		ListName:   listName,
		GroupName:  groupName,
		GroupRank:  groupRank,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x16,
		Read: ReadExtAddPlayerName,
		Size: ExtAddPlayerNameSize,
		Type: packets.ServerOnly,
		Name: "Ext Add Player Name (CPE)",
	})
}
