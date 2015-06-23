//go:generate protocol_generator -file=$GOFILE -direction=serverbound -state=play -package=modern

package modern

import "github.com/kurafuto/kyubu/packets"

// Packet ID: 0x01
type ClientMessage struct {
	Message string
}

// Packet ID: 0x14
type ClientTabComplete struct {
	Text          string
	HasPosition   bool
	LookedAtBlock packets.Position `if:".HasPosition"`
}

// Packet ID: 0x16
type ClientStatus struct {
	ActionID packets.VarInt
}

// Packet ID: 0x17
type ClientPluginMessage struct {
	Channel string
	Data    []byte `length:"rest"`
}
