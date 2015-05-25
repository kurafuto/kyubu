// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=test.go -direction=anomalous -state=handshake -package=main

package main

import (
	"encoding/binary"
	"github.com/sysr-q/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Handshake, packets.Anomalous, 0xfe, func() packets.Packet { return &TestPacket{} })
}

func (t *TestPacket) Id() byte {
	return 0xfe // 254
}

func (t *TestPacket) Encode(ww io.Writer) (err error) {
	// Encoding: b (bool)
	tmp0 := byte(0)
	if t.b {
		tmp0 = byte(1)
	}
	if err = binary.Write(ww, binary.BigEndian, tmp0); err != nil {
		return
	}

	// Encoding: i8 (int8)
	if err = binary.Write(ww, binary.BigEndian, t.i8); err != nil {
		return
	}

	// Encoding: ui8 (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.ui8); err != nil {
		return
	}

	// Encoding: i16 (int16)
	if err = binary.Write(ww, binary.BigEndian, t.i16); err != nil {
		return
	}

	// Encoding: ui16 (uint16)
	if err = binary.Write(ww, binary.BigEndian, t.ui16); err != nil {
		return
	}

	// Encoding: i32 (int32)
	if err = binary.Write(ww, binary.BigEndian, t.i32); err != nil {
		return
	}

	// Encoding: i64 (int64)
	if err = binary.Write(ww, binary.BigEndian, t.i64); err != nil {
		return
	}

	// Encoding: f32 (float32)
	if err = binary.Write(ww, binary.BigEndian, t.f32); err != nil {
		return
	}

	// Encoding: f64 (float64)
	if err = binary.Write(ww, binary.BigEndian, t.f64); err != nil {
		return
	}

	// Encoding: V32 (packets.VarInt)
	tmp1 := make([]byte, binary.MaxVarintLen64)
	tmp2 := packets.PutVarint(tmp1, int64(t.V32))
	if err = binary.Write(ww, binary.BigEndian, tmp1[:tmp2]); err != nil {
		return
	}

	// Encoding: V64 (packets.VarLong)
	tmp3 := make([]byte, binary.MaxVarintLen64)
	tmp4 := packets.PutVarint(tmp3, int64(t.V64))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp4]); err != nil {
		return
	}

	// Encoding: C (packets.Chunk)

	// Encoding: M (packets.Metadata)

	// Encoding: S (packets.Slot)

	// Encoding: O (packets.ObjectData)

	// Encoding: N (packets.NBT)

	// Encoding: P (packets.Position)
	if err = binary.Write(ww, binary.BigEndian, t.P); err != nil {
		return
	}

	// Encoding: A (packets.Angle)
	if err = binary.Write(ww, binary.BigEndian, t.A); err != nil {
		return
	}

	// Encoding: U (packets.UUID)
	if err = binary.Write(ww, binary.BigEndian, t.U[:]); err != nil {
		return
	}

	return
}

func (t *TestPacket) Decode(rr io.Reader) (err error) {
	// Decoding: b (bool)
	var tmp0 [1]byte
	if _, err = rr.Read(tmp0[:1]); err != nil {
		return
	}
	t.b = tmp0[0] == 0x01

	// Decoding: i8 (int8)
	if err = binary.Read(rr, binary.BigEndian, t.i8); err != nil {
		return
	}

	// Decoding: ui8 (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.ui8); err != nil {
		return
	}

	// Decoding: i16 (int16)
	if err = binary.Read(rr, binary.BigEndian, t.i16); err != nil {
		return
	}

	// Decoding: ui16 (uint16)
	if err = binary.Read(rr, binary.BigEndian, t.ui16); err != nil {
		return
	}

	// Decoding: i32 (int32)
	if err = binary.Read(rr, binary.BigEndian, t.i32); err != nil {
		return
	}

	// Decoding: i64 (int64)
	if err = binary.Read(rr, binary.BigEndian, t.i64); err != nil {
		return
	}

	// Decoding: f32 (float32)
	if err = binary.Read(rr, binary.BigEndian, t.f32); err != nil {
		return
	}

	// Decoding: f64 (float64)
	if err = binary.Read(rr, binary.BigEndian, t.f64); err != nil {
		return
	}

	// Decoding: V32 (packets.VarInt)
	tmp1, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.V32 = packets.VarInt(tmp1)

	// Decoding: V64 (packets.VarLong)
	tmp2, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.V64 = packets.VarLong(tmp2)

	// Decoding: C (packets.Chunk)

	// Decoding: M (packets.Metadata)

	// Decoding: S (packets.Slot)

	// Decoding: O (packets.ObjectData)

	// Decoding: N (packets.NBT)

	// Decoding: P (packets.Position)
	if err = binary.Read(rr, binary.BigEndian, t.P); err != nil {
		return
	}

	// Decoding: A (packets.Angle)
	if err = binary.Read(rr, binary.BigEndian, t.A); err != nil {
		return
	}

	// Decoding: U (packets.UUID)
	if err = binary.Read(rr, binary.BigEndian, t.U); err != nil {
		return
	}

	return
}
