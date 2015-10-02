package packets

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrNegativeLength = errors.New("kyubu: packet has negative length")
	ErrUnknownPacket  = errors.New("kyubu: packet has unknown id")
	ErrLostSync       = errors.New("kyubu: lost sync on packets")
	ErrNotCompressed  = errors.New("kyubu: packet not compressed correctly")
)

type Parser struct {
	r    io.Reader
	w    io.Writer
	conn net.Conn

	zlibReader io.ReadCloser
	zlibWriter *zlib.Writer // TODO

	State     State
	Direction PacketDirection

	compressionThreshold int
	timeout              time.Duration
}

func NewParser(c net.Conn, state State, dir PacketDirection) *Parser {
	return &Parser{
		r:    c,
		w:    c,
		conn: c,

		State:     state,
		Direction: dir,

		compressionThreshold: -1,
		timeout:              15 * time.Second,
	}
}

func (p *Parser) SetCompression(threshold int) {
	p.compressionThreshold = threshold
}

// Recv returns the next packet parsed out of the stream.
func (p *Parser) Recv() (Packet, error) {
	p.conn.SetReadDeadline(time.Now().Add(p.timeout))

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

	var r *bytes.Reader
	r = bytes.NewReader(b)

	// If compression is enabled, we need to know how much data to expect first.
	if p.compressionThreshold >= 0 {
		// newSize is the decompressed size of the packet.
		newSize, err := ReadVarint(r)
		if err != nil {
			return nil, err
		}

		/*
			// Technically, this is meant to disconnect you, but we can just
			// ignore this case. Most parsers seem to.
			if newSize < p.compressionThreshold {
				return nil, ErrNotCompressed
			}
		*/

		// If newSize == 0, it's not compressed so we can just carry on anyway.
		if newSize != 0 {
			if p.zlibReader == nil {
				// If we don't already have a zlib reader, set one up.
				p.zlibReader, err = zlib.NewReader(r)
				if err != nil {
					return nil, err
				}
			} else {
				// If we do, reset with the current packet.
				// NOTE: This might not be efficient. Will have to test.
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

	// f() returns an empty Packet, that we'll decode into.
	f, ok := Packets[p.State][p.Direction.Flip()][byte(id)]
	if !ok {
		//fmt.Errorf("kyubu: unknown packet %s:%#.2x", p.State, id)
		return nil, ErrUnknownPacket
	}

	packet := f()

	// XXX: Can we pool packet structs, or is that not feasible since there are
	// underlying fields we don't know about in the generic type?
	if err := packet.Decode(r); err != nil {
		return nil, err
	}

	// If there are bytes left over, the decoding step missed something.
	if r.Len() > 0 {
		//fmt.Errorf("kyubu: Lost sync on packet %s:%#.2x, %d bytes left", p.State, id, r.Len())
		return nil, ErrLostSync
	}

	return packet, nil
}

// Send encodes (and compresses, if required) a packet, and sends it.
func (p *Parser) Send(packet Packet) error {
	p.conn.SetWriteDeadline(time.Now().Add(p.timeout))
	// TODO
	return nil
}

func (p *Parser) EnableEncryption(key []byte) error {
	// TODO
	return nil
}
