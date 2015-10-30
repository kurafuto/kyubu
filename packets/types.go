// types.go has type aliases for all of the protocol types documented in:
// http://wiki.vg/Protocol#Data_types
// NOTE: We only alias non Go-builtin types. bool, int8, etc aren't here.

package packets

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/kurafuto/kyubu/format"
)

var Endianness = binary.BigEndian

// (>= 1, <= 5) -2147483648 to 2147483647
type VarInt int32

// (>= 1, <= 10) -9223372036854775808 to 9223372036854775807
type VarLong int64

type Chat format.AnyComponent

/* // TODO: Implement these placeholders.
type Chunk bool
type Metadata bool
type Slot bool
type ObjectData bool
type NBT bool
*/

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

type UUID [16]byte

func ReadBool(r io.Reader) (bool, error) {
	var b [1]byte
	if _, err := r.Read(b[:1]); err != nil {
		return false, err
	}
	return b[0] == 0x01, nil
}

func WriteBool(w io.Writer, b bool) error {
	var bb byte = 0
	if b {
		bb = 0x01
	}
	return binary.Write(w, Endianness, bb)
}

func ReadString(r io.Reader) (string, error) {
	n, err := ReadVarint(r)
	if err != nil {
		return "", err
	}

	b := make([]byte, n)
	x, err := r.Read(b)

	if err != nil {
		return "", err
	} else if int64(x) != n { // TODO: Needed?
		return "", errors.New("didn't read enough bytes for string")
	}

	return string(b), nil
}

func WriteString(w io.Writer, s string) error {
	x := make([]byte, binary.MaxVarintLen32)
	b := []byte(s)
	n := PutVarint(x, int64(len(b)))
	if err := binary.Write(w, Endianness, x[:n]); err != nil {
		return err
	}
	if err := binary.Write(w, Endianness, b); err != nil {
		return err
	}

	return nil
}

func WriteVarint(w io.Writer, v VarInt) error {
	x := make([]byte, binary.MaxVarintLen32)
	n := PutVarint(x, int64(v)) // XXX: Should we use VarLong instead?

	if err := binary.Write(w, Endianness, x[:n]); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////
// NOTE: Copied + modified from Go stdlib source code.
// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
///////////////////////////////////////////////////////

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

// Uvarint decodes a uint64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 64 bits (overflow)
//              and -n is the number of bytes read
//
func Uvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, b := range buf {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, -(i + 1) // overflow
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

// PutVarint encodes an int64 into buf and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.
func PutVarint(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutUvarint(buf, ux)
}

// Varint decodes an int64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 with the following meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 64 bits (overflow)
//              and -n is the number of bytes read
//
func Varint(buf []byte) (int64, int) {
	ux, n := Uvarint(buf) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n
}

var overflow = errors.New("binary: varint overflows a 64-bit integer")

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func ReadUvarint(r io.Reader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		b, err := readByte(r) //r.ReadByte()
		if err != nil {
			return x, err
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, overflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

// ReadVarint reads an encoded signed integer from r and returns it as an int64.
func ReadVarint(r io.Reader) (int64, error) {
	ux, err := ReadUvarint(r) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, err
}

func readByte(r io.Reader) (byte, error) {
	if rr, ok := r.(io.ByteReader); ok {
		return rr.ReadByte()
	}
	var b [1]byte
	_, err := r.Read(b[:1])
	return b[0], err
}
