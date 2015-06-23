//go:generate protocol_generator -file=$GOFILE -direction=serverbound -state=status -package=modern

package modern

// Packet ID: 0x00
type StatusRequest struct {
}

// Packet ID: 0x01
type StatusPing struct {
	Time int64
}
