// types.go has type aliases for all of the protocol types documented in:
// http://wiki.vg/Protocol#Data_types
// NOTE: We only alias non Go-builtin types. bool, int8, etc aren't here.

package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var Endianness = binary.BigEndian

/* Eventually we'll just need these:
//type VarInt int32
//type VarLong int64

type Chunk x
type Metadata x
type Slot x
type ObjectData x
type NBT x
//type Position uint64
//type Angle x
type UUID x
*/

// (>= 1, <= 5) -2147483648 to 2147483647
type VarInt int32

// (>= 1, <= 10) -9223372036854775808 to 9223372036854775807
type VarLong int64

// (8) x (-33554432 to 33554431), y (-2048 to 2047), z (-33554432 to 33554431)
// x as a 26-bit integer, followed by y as a 12-bit integer, followed by z as a 26-bit integer
// Code lifted from thinkofdeath/steven:protocol/handshaking.go
type Position uint64

func NewPosition(x, y, z int) Position {
	return ((Position(x) & 0x3FFFFFF) << 38) | ((Position(y) & 0xFFF) << 26) | (Position(z) & 0x3FFFFFF)
}

func (p Position) X() int {
	return int(int64(p) >> 38)
}

func (p Position) Y() int {
	return int((int64(p) >> 26) & 0xFFF)
}

func (p Position) Z() int {
	return int(int64(p) << 38 >> 38)
}

func (p Position) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X(), p.Y(), p.Z())
}

// (1) Rotation angle in steps of 1/256 of a full turn
// Unsigned to get 0-255, rather than -128 to 127
type Angle uint8

//////// TODO: Remove below

// (>= 1, <= 2147483652) A sequence of Unicode code points
// UTF-8 string prefixed with its size in bytes as a VarInt
type String string

func (t *String) Encode(buf *bytes.Buffer) (err error) {
	// TODO: defer + recover() for PutVarint.
	// TODO: This might not be efficient. Pool []byte arrays?
	x := make([]byte, binary.MaxVarintLen32)
	b := []byte(*t)

	n := binary.PutVarint(x, int64(len(b)))

	err = binary.Write(buf, Endianness, x[:n])
	if err != nil {
		return
	}

	err = binary.Write(buf, Endianness, b)
	if err != nil {
		return
	}
	return nil
}

func (t *String) Decode(buf *bytes.Buffer) error {
	n, err := binary.ReadVarint(buf)
	if err != nil {
		return err
	}

	b := make([]byte, n)
	n1, err := buf.Read(b)
	if err != nil {
		return err
	} else if int64(n1) != n {
		// TODO: delta := n - n1; Read()
		return errors.New("didn't read enough bytes for string")
	}

	*t = String(b)
	return nil
}
