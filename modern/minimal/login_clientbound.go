//go:generate protocol_generator -file=$GOFILE -direction=clientbound -state=login -package=modern

package modern

import "github.com/kurafuto/kyubu/packets"

// Packet ID: 0x00
type LoginDisconnect struct {
	Reason packets.Chat `as:"json"`
}

// Packet ID: 0x01
type EncryptionRequest struct {
	ServerID    string
	PublicKey   []byte `length:"packets.VarInt"`
	VerifyToken []byte `length:"packets.VarInt"`
}

// Packet ID: 0x02
type LoginSuccess struct {
	UUID     string
	Username string
}

// Packet ID: 0x03
type SetInitialCompression struct {
	Threshold packets.VarInt
}
