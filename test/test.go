//go:generate protocol_generator -file=$GOFILE -direction=anomalous -state=handshake -package=main
package main

import "github.com/kurafuto/kyubu/packets"

// This packet contains all the types that protocol_generator can work with.
// It basically exists to make sure the output is type safe and compilable.
//
// Packet ID: 0xfe
type TestPacket struct {
	// The built-in types.
	b    bool
	i8   int8
	ui8  uint8
	i16  int16
	ui16 uint16
	i32  int32
	i64  int64
	f32  float32
	f64  float64

	// Custom types
	V32 packets.VarInt
	V64 packets.VarLong
	CH  packets.Chat
	P   packets.Position
	A   packets.Angle
	U   packets.UUID

	// These aren't implemented, probably gonna error.
	/*
		C packets.Chunk
		M packets.Metadata
		S packets.Slot
		O packets.ObjectData
		N packets.NBT
	*/
}

// This packet just tests nesting structs for "complex" (eh..) packets.
// Not really required now, but it's implemented in the generator.
// Should end up like: t.T.X = ...
//
// Packet ID: 0xfd
type Test struct {
	T struct {
		X uint16
		Y string
		Z packets.VarInt
	}
	TXs []TX `length:"packets.VarInt"`
}

type TX struct {
	X uint16
	Y string
	Z packets.VarInt
}
