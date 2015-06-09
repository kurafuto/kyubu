//go:generate protocol_generator -file=$GOFILE -direction=anomalous -state=handshake -package=main
package main

// Packet ID: 0x42
type X struct {
	Ys []Y `length:"packets.VarInt"`
}

type Y struct {
	b bool
}
