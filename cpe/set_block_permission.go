package cpe

import "github.com/sysr-q/kyubu/packets"

const SetBlockPermissionSize = (packets.ByteSize + packets.ByteSize + packets.ByteSize + packets.ByteSize)

type SetBlockPermission struct {
	PacketId       byte
	BlockType      byte
	AllowPlacement byte
	AllowDeletion  byte
}

func (p SetBlockPermission) Id() byte {
	return p.PacketId
}

func (p SetBlockPermission) Size() int {
	return SetBlockPermissionSize
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
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x1c,
		Read: ReadSetBlockPermission,
		Size: SetBlockPermissionSize,
		Type: packets.ServerOnly,
		Name: "Set Block Permission (CPE)",
	})
}
