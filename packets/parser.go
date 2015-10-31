package packets

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
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
	zlibWriter *zlib.Writer

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

	buf := &bytes.Buffer{}

	// Write the packet ID, then encode it into the buffer.
	if err := WriteVarint(buf, VarInt(packet.Id())); err != nil {
		return err
	}
	if err := packet.Encode(buf); err != nil {
		return err
	}

	uncompSize := 0
	extra := 0 // Just in case we have to compress the packet.

	// We only need to compress the packet if it's above the threshold.
	if p.compressionThreshold >= 0 && buf.Len() > p.compressionThreshold {
		// We're compressing the packet into this buffer.
		cBuf := &bytes.Buffer{}

		if p.zlibWriter == nil {
			// If we don't already have a zlib writer, set one up.
			p.zlibWriter, _ = zlib.NewWriterLevel(cBuf, zlib.BestSpeed)
		} else {
			p.zlibWriter.Reset(cBuf)
		}

		// Store the size it should end up being on the receiving end.
		uncompSize = buf.Len()
		extra = binary.Size(uncompSize)

		if _, err := buf.WriteTo(p.zlibWriter); err != nil {
			return err
		}

		// We have to close the zlib writer, otherwise it won't flush and
		// actually compress the data.
		if err := p.zlibWriter.Close(); err != nil {
			return err
		}

		buf = cBuf
	}

	if err := WriteVarint(p.w, VarInt(buf.Len()+extra)); err != nil {
		return err
	}

	if p.compressionThreshold >= 0 {
		// If compression is on, tell them how big it should end up being.
		if err := WriteVarint(p.w, VarInt(uncompSize)); err != nil {
			return err
		}
	}

	// Finally, the rest of the packet buffer.
	_, err := buf.WriteTo(p.w)
	return err
}

func (p *Parser) EnableEncryption(key []byte) error {
	// TODO: Set up EnableEncryption() with CFB8.
	// http://wiki.vg/Protocol_Encryption, http://git.io/vlZiG
	return nil
}
