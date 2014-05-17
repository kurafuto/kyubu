package packets

import (
	"fmt"
	"io"
)

type Parser interface {
	Next() (Packet, error)
}

type parser struct {
	r io.Reader
	//C <-chan Packet
}

func (p *parser) Next() (Packet, error) {
	// Parse out a single byte id
	idByte := make([]byte, 1)
	n, err := p.r.Read(idByte)
	if err != nil && err == io.EOF {
		return nil, err
	} else if err != nil {
		return nil, err
	} else if n != len(idByte) {
		return nil, fmt.Errorf("kyubu: Read failed for id, wanted %d bytes, got %d", len(idByte), n)
	}

	// Packet ID
	id := idByte[0]

	// Figure out what packet we're expecting.
	packetInfo, ok := Packets[id]
	if !ok {
		return nil, fmt.Errorf("kyubu: Unrecognized packet id %#.2x", id)
	}

	packetData := make([]byte, packetInfo.Size)
	if packetInfo.Size == 1 {
		// We only wanted the packet id.
		packetData = idByte
	} else {
		// We've already read 1 byte for the packet id
		restSize := packetInfo.Size - ByteSize
		restData := make([]byte, restSize)
		n, err := p.r.Read(restData)
		if err != nil && err == io.EOF {
			return nil, fmt.Errorf("kyubu: EOF reading packet %#.2x, got %d", id, n)
		} else if err != nil {
			return nil, err
		} else if n != restSize {
			// We still want `delta` bytes. Let's try to get them.
			delta := restSize - n
			rest := make([]byte, delta)

			an, err := p.r.Read(rest)
			if err != nil {
				return nil, fmt.Errorf("kyubu: Error reading packet %#.2x: %s", id, err)
			}
			// Slice so we don't accidentally overflow the packet size.
			restData = append(restData[:n], rest...)
			n += an
		}

		if n != restSize {
			return nil, fmt.Errorf("kyubu: got %d bytes, wanted %d, reading packet %#.2x", n+1, packetInfo.Size, id)
		}

		// Ensure we've also got the packet id.
		packetData = append([]byte{id}, restData...)
	}

	packet, err := packetInfo.Read(packetData)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func NewParser(r io.Reader) Parser {
	return &parser{
		r: r,
		//C: make(<-chan Packet),
	}
}
