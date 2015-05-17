package cpe

import "github.com/sysr-q/kyubu/packets"

type SetBlockPermission struct {
	PacketId       byte
	BlockType      byte
	AllowPlacement byte
	AllowDeletion  byte
}

func (p SetBlockPermission) Id() byte {
	return 0x1c
}

func (p SetBlockPermission) Size() int {
	return packets.ReflectSize(p)
}

func (p SetBlockPermission) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p SetBlockPermission) String() string {
	return "BlockPermissions"
}

func ReadSetBlockPermission(b []byte) (packets.Packet, error) {
	var p SetBlockPermission
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewSetBlockPermission(blockType, allowPlacement, allowDeletion byte) (p *SetBlockPermission, err error) {
	p = &SetBlockPermission{
		PacketId:       0x1c,
		BlockType:      blockType,
		AllowPlacement: allowPlacement,
		AllowDeletion:  allowDeletion,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x1c,
		Read:      ReadSetBlockPermission,
		Size:      packets.ReflectSize(&SetBlockPermission{}),
		Direction: packets.Anomalous,
		Name:      "Set Block Permission (CPE)",
	})
}
