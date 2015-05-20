// types.go has type aliases for all of the protocol types documented in:
// http://wiki.vg/Protocol#Data_types

package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var Endianness = binary.BigEndian

// (1) false, true
// Value can be either true (0x01) or false (0x00)
type Boolean bool

func (t *Boolean) Encode(buf *bytes.Buffer) error {
	if *t {
		return binary.Write(buf, Endianness, byte(1))
	} else {
		return binary.Write(buf, Endianness, byte(0))
	}
}
func (t *Boolean) Decode(buf *bytes.Buffer) (err error) {
	b := make([]byte, 1)
	_, err = buf.Read(b)
	if err != nil {
		return
	}
	if b[0] == 0x01 {
		*t = Boolean(true)
	} else if b[0] == 0x00 {
		*t = Boolean(false)
	} else {
		err = fmt.Errorf("neither 0x00 or 0x01: % x", b[0])
		return
	}
	return
}

// (1) -128 to 127
// Signed 8-bit integer, two's complement
type Byte int8

func (t *Byte) Encode(buf *bytes.Buffer) error {
	return binary.Write(buf, Endianness, t)
}

func (t *Byte) Decode(buf *bytes.Buffer) (err error) {
	i := make([]byte, 1)
	_, err = buf.Read(i)
	if err != nil {
		return
	}
	*t = Byte(i[0])
	return
}

// (1) 0 to 255
// Unsigned 8-bit integer
type UByte uint8

func (t *UByte) Encode(buf *bytes.Buffer) error {
	return binary.Write(buf, Endianness, t)
}

func (t *UByte) Decode(buf *bytes.Buffer) (err error) {
	i := make([]byte, 1)
	_, err = buf.Read(i)
	if err != nil {
		return
	}
	*t = UByte(i[0])
	return
}

// (2) -32768 to 32767
// Signed 16-bit integer, two's complement
type Short int16

func (t *Short) Encode(buf *bytes.Buffer) error {
	return binary.Write(buf, Endianness, t)
}

func (t *Short) Decode(buf *bytes.Buffer) (err error) {
	err = binary.Read(buf, Endianness, t)
	return
}

// (2) 0 to 65535
// Unsigned 16-bit integer
type UShort uint16

// (4) -2147483648 to 2147483647
// Signed 32-bit integer, two's complement
type Int int32

// (8) -9223372036854775808 to 9223372036854775807
// Signed 64-bit integer, two's complement
type Long int64

// (4) Single-precision 32-bit IEEE 754 floating point
type Float float32

// (8) Double-precision 64-bit IEEE 754 floating point
type Double float64

// (>= 1, <= 2147483652) A sequence of Unicode code points
// UTF-8 string prefixed with its size in bytes as a VarInt
type String string

// (>= 1, <= 5) -2147483648 to 2147483647
// Read with encoding/binary.ReadVarint
// TODO: Cast it to/from an int32
type VarInt int32

// (>= 1, <= 10) -9223372036854775808 to 9223372036854775807
// Read with encoding/binary.ReadVarint
type VarLong int64

// TODO: Implement
/*
type Chunk x
type Metadata x
type Slot x
type ObjectData x
type NBTTag x
*/

// (8) x (-33554432 to 33554431), y (-2048 to 2047), z (-33554432 to 33554431)
// x as a 26-bit integer, followed by y as a 12-bit integer, followed by z as a 26-bit integer
type Position int64

// (1) Rotation angle in steps of 1/256 of a full turn
type Angle byte

// TODO: Implement
/*
type UUID x
*/
