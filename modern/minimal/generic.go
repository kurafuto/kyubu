package minimal

import (
	"encoding/binary"
	"io"
)

func init() {
	// TODO: Register the packets not implemented in other files.
}

type GenericPacket struct {
	Data []byte
}

func (t *GenericPacket) Encode(ww io.Writer) error {
	return binary.Write(ww, binary.BigEndian, t.Data)
}

func (t *GenericPacket) Decode(rr io.Reader) error {
	b, err := io.ReadAll(rr)
	if err != nil {
		return err
	}

	t.Data = b
	return nil
}
