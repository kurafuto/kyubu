// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=status_clientbound.go -direction=serverbound -state=status -package=minimal

package minimal

import (
	"encoding/binary"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Status, packets.ServerBound, 0x00, func() packets.Packet { return &StatusResponse{} })
	packets.Register(packets.Status, packets.ServerBound, 0x01, func() packets.Packet { return &StatusPong{} })
}

func (t *StatusResponse) Id() byte {
	return 0x00 // 0
}

func (t *StatusResponse) Encode(ww io.Writer) (err error) {
	// Encoding: Status (StatusReply)
	// Unable to encode: t.Status (StatusReply)

	return
}

func (t *StatusResponse) Decode(rr io.Reader) (err error) {
	// Decoding: Status (StatusReply)
	// Unable to decode: t.Status (StatusReply)

	return
}

func (t *StatusPong) Id() byte {
	return 0x01 // 1
}

func (t *StatusPong) Encode(ww io.Writer) (err error) {
	// Encoding: Time (int64)
	if err = binary.Write(ww, binary.BigEndian, t.Time); err != nil {
		return
	}

	return
}

func (t *StatusPong) Decode(rr io.Reader) (err error) {
	// Decoding: Time (int64)
	if err = binary.Read(rr, binary.BigEndian, t.Time); err != nil {
		return
	}

	return
}
