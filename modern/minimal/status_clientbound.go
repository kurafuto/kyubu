//go:generate protocol_generator -file=$GOFILE -direction=serverbound -state=status -package=minimal

package minimal

// Packet ID: 0x00
type StatusResponse struct {
	Status StatusReply `as:"json"`
}

// Packet ID: 0x01
type StatusPong struct {
	Time int64
}
