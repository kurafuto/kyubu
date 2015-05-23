// Package packets contains the necessary magic to create, parse, serialize and
// register Minecraft packets.
package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// Size constants from spec, used to deduce packet size.
const (
	ByteSize   = 1
	SByteSize  = 1
	ShortSize  = 2
	IntSize    = 4
	StringSize = 64
	BytesSize  = 1024
)

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

/*
// TODO: Remove below
var (
	// ServerPackets is a map of packet ids to serverbound packets.
	ServerPackets = map[byte]*PacketInfo{}
	// ClientPackets is a map of packet ids to clientbound packets.
	ClientPackets = map[byte]*PacketInfo{}

	// AnomalousPackets is a map of packet ids to "odd" packets. Basically
	// packets which are declared as anomalous, which is kept around for Classic
	// stuff (for now).
	AnomalousPackets = map[byte]*PacketInfo{}
)

type PacketInfo struct {
	Id   byte
	Name string // Human presentable

	Read      readFunc
	Size      int
	Direction PacketDirection
}

type readFunc func([]byte) (Packet, error)

// Register allows users to register additional packets with the internal parser.
// These packets will be recognised, parsed and sent to the parser channel.
func Register(p *PacketInfo) {
	if p.Direction == ServerBound {
		ServerPackets[p.Id] = p
	} else if p.Direction == ClientBound {
		ClientPackets[p.Id] = p
	} else if p.Direction == Anomalous {
		AnomalousPackets[p.Id] = p
	}
}
*/

// PacketWrapper wraps an array of bytes into a bytes.Buffer, and exposes
// various Read/Write type helpers. It also optionally (and by default) strips
// padding characters, as defined by the Minecraft Classic Protocol spec.
type PacketWrapper struct {
	Buffer *bytes.Buffer
	Strip  bool
}

func (p *PacketWrapper) WriteByte(b byte) error {
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) WriteSByte(i int8) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteShort(i int16) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteInt(i int32) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteString(s string) error {
	if len(s) > StringSize {
		return fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
	}
	// Pad with spaces, up to 64 bytes.
	s += strings.Repeat(" ", StringSize-len(s))
	b := []byte(s)
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) WriteBytes(b []byte) error {
	if len(b) > BytesSize {
		return fmt.Errorf("kyubu/packets: cannot write over %d bytes", BytesSize)
	}
	// Pad with 0x00, up to 1024 bytes.
	b = append(b, bytes.Repeat([]byte{0}, BytesSize-len(b))...)
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) ReadByte() (i byte, err error) {
	b := make([]byte, 1)
	_, err = p.Buffer.Read(b)
	if err != nil {
		return
	}
	i = b[0]
	return
}

func (p *PacketWrapper) ReadSByte() (i int8, err error) {
	err = binary.Read(p.Buffer, Endianness, &i)
	return
}

func (p *PacketWrapper) ReadShort() (i int16, err error) {
	err = binary.Read(p.Buffer, Endianness, &i)
	return
}

func (p *PacketWrapper) ReadInt() (i int32, err error) {
	err = binary.Read(p.Buffer, Endianness, &i)
	return
}

func (p *PacketWrapper) ReadString() (s string, err error) {
	// Strings are 64 bytes, padded with spaces (0x20).
	b := make([]byte, StringSize)
	_, err = p.Buffer.Read(b)
	if err != nil {
		return
	}
	s = string(b)
	if p.Strip {
		s = strings.TrimRight(s, " ")
	}
	return
}

func (p *PacketWrapper) ReadBytes() (b []byte, err error) {
	b = p.Buffer.Next(BytesSize)
	if p.Strip {
		b = bytes.TrimRight(b, string([]byte{0}))
	}
	return
}

func NewPacketWrapper(b []byte) (p *PacketWrapper) {
	p = &PacketWrapper{
		Buffer: bytes.NewBuffer(b),
		Strip:  true,
	}
	return
}
