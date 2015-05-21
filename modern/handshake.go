//go:generate protocol_generator -file=$GOFILE -direction=serverbound -state=handshake

package modern

import "github.com/sysr-q/kyubu/packets"

// Packet ID: 0x00
type Handshake struct {
	ProtocolVersion packets.VarInt
	Address         string
	Port            uint16
	NextState       packets.VarInt
}
