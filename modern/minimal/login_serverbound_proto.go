// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=login_serverbound.go -direction=serverbound -state=login -package=minimal

package minimal

import (
	"encoding/binary"
	"errors"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Login, packets.ServerBound, 0x00, func() packets.Packet { return &LoginStart{} })
	packets.Register(packets.Login, packets.ServerBound, 0x01, func() packets.Packet { return &EncryptionResponse{} })
}

func (t *LoginStart) Id() byte {
	return 0x00 // 0
}

func (t *LoginStart) Encode(ww io.Writer) (err error) {
	// Encoding: Username (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Username)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	return
}

func (t *LoginStart) Decode(rr io.Reader) (err error) {
	// Decoding: Username (string)
	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp1 := make([]byte, tmp0)
	tmp2, err := rr.Read(tmp1)
	if err != nil {
		return
	} else if int64(tmp2) != tmp0 {
		return errors.New("didn't read enough bytes for string")
	}
	t.Username = string(tmp1)

	return
}

func (t *EncryptionResponse) Id() byte {
	return 0x01 // 1
}

func (t *EncryptionResponse) Encode(ww io.Writer) (err error) {
	// Encoding: SharedSecret ()
	// Unable to encode: t.SharedSecret ()

	// Encoding: VerifyToken ()
	// Unable to encode: t.VerifyToken ()

	return
}

func (t *EncryptionResponse) Decode(rr io.Reader) (err error) {
	// Decoding: SharedSecret ()
	// Unable to decode: t.SharedSecret ()

	// Decoding: VerifyToken ()
	// Unable to decode: t.VerifyToken ()

	return
}
