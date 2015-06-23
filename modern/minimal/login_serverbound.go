//go:generate protocol_generator -file=$GOFILE -direction=serverbound -state=login -package=modern

package modern

// Packet ID: 0x00
type LoginStart struct {
	Username string
}

// Packet ID: 0x01
type EncryptionResponse struct {
	SharedSecret []byte `length:"packets.VarInt"`
	VerifyToken  []byte `length:"packets.VarInt"`
}
