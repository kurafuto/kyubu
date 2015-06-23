package packets

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNegativeLength = errors.New("kyubu: packet has negative length")
)

type Parser struct {
	r io.Reader

	State     State
	Direction PacketDirection

	CompressionThreshold int

	zlibReader io.ReadCloser
}

func (p *Parser) Next() (Packet, error) {
	length, err := ReadVarint(p.r)
	if err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, ErrNegativeLength
	}

	b := make([]byte, length)
	if _, err := io.ReadFull(p.r, b); err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)

	if p.CompressionThreshold >= 0 {
		newSize, err := ReadVarint(r)
		if err != nil {
			return nil, err
		}

		if newSize != 0 {
			if p.zlibReader == nil {
				p.zlibReader, err = zlib.NewReader(r)
				if err != nil {
					return nil, err
				}
			} else {
				err = p.zlibReader.(zlib.Resetter).Reset(r, nil)
				if err != nil {
					return nil, err
				}
			}

			data := make([]byte, newSize)
			_, err := io.ReadFull(p.zlibReader, data)
			if err != nil {
				return nil, err
			}
			r = bytes.NewReader(data)
		}
	}

	id, err := ReadVarint(r)
	if err != nil {
		return nil, err
	}

	f, ok := Packets[p.State][p.Direction.Flip()][byte(id)]
	if !ok {
		return nil, fmt.Errorf("kyubu: unknown packet %s:%#.2x", p.State, id)
	}

	packet := f()

	if err := packet.Decode(r); err != nil {
		return nil, err
	}

	if r.Len() > 0 {
		return nil, fmt.Errorf("kyubu: Lost sync on packet %s:%#.2x, %d bytes left", p.State, id, r.Len())
	}

	return packet, nil
}

func NewParser(r io.Reader, state State, dir PacketDirection) *Parser {
	return &Parser{
		r: r,

		State:     state,
		Direction: dir,

		CompressionThreshold: -1,
	}
}
