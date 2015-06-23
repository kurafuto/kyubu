// Package packets contains the necessary magic to create, parse, serialize and
// register Minecraft packets.
package packets

import "io"

// PacketDirection lets us define whether packets are purely C->S or S->C.
// This will change what the parser reads, depending on its direction.
type PacketDirection int

func (p *PacketDirection) Flip() PacketDirection {
	if *p == ServerBound {
		return ClientBound
	} else if *p == ClientBound {
		return ServerBound
	} else {
		return Anomalous
	}
}

// ServerBound and ClientBound let us designate a parser's direction. This
// defines what kind of packets it'll parse/serialize.
const (
	ServerBound PacketDirection = iota
	ClientBound
	Anomalous
)

type State int

const (
	Handshake State = iota
	Status
	Login
	Play
)

type Packet interface {
	Id() byte
	Encode(io.Writer) error
	Decode(io.Reader) error
}

// TODO: This crap.
type PacketFunc func() Packet

var Packets [4][3]map[byte]PacketFunc

func Register(state State, dir PacketDirection, id byte, f PacketFunc) {
	Packets[state][dir][id] = f
}
