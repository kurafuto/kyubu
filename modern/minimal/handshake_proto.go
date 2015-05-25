// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=handshake.go -direction=serverbound -state=handshake -package=minimal

package minimal

import (
	"encoding/binary"
	"errors"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Handshake, packets.ServerBound, 0x00, func() packets.Packet { return &Handshake{} })
}

func (t *Handshake) Id() byte {
	return 0x00 // 0
}

func (t *Handshake) Encode(ww io.Writer) (err error) {
	// Encoding: ProtocolVersion (packets.VarInt)
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(t.ProtocolVersion))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return
	}

	// Encoding: Address (string)
	tmp2 := make([]byte, binary.MaxVarintLen32)
	tmp3 := []byte(t.Address)
	tmp4 := packets.PutVarint(tmp2, int64(len(tmp3)))
	if err = binary.Write(ww, binary.BigEndian, tmp2[:tmp4]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp3); err != nil {
		return
	}

	// Encoding: Port (uint16)
	if err = binary.Write(ww, binary.BigEndian, t.Port); err != nil {
		return
	}

	// Encoding: NextState (packets.VarInt)
	tmp5 := make([]byte, binary.MaxVarintLen64)
	tmp6 := packets.PutVarint(tmp5, int64(t.NextState))
	if err = binary.Write(ww, binary.BigEndian, tmp5[:tmp6]); err != nil {
		return
	}

	return
}

func (t *Handshake) Decode(rr io.Reader) (err error) {
	// Decoding: ProtocolVersion (packets.VarInt)
	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.ProtocolVersion = packets.VarInt(tmp0)

	// Decoding: Address (string)
	tmp1, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp2 := make([]byte, tmp1)
	tmp3, err := rr.Read(tmp2)
	if err != nil {
		return
	} else if int64(tmp3) != tmp1 {
		return errors.New("didn't read enough bytes for string")
	}
	t.Address = string(tmp2)

	// Decoding: Port (uint16)
	if err = binary.Read(rr, binary.BigEndian, t.Port); err != nil {
		return
	}

	// Decoding: NextState (packets.VarInt)
	tmp4, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.NextState = packets.VarInt(tmp4)

	return
}
