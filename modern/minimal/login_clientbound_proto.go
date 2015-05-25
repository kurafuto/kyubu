// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=login_clientbound.go -direction=clientbound -state=login -package=minimal

package minimal

import (
	"encoding/binary"
	"errors"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Login, packets.ClientBound, 0x00, func() packets.Packet { return &LoginDisconnect{} })
	packets.Register(packets.Login, packets.ClientBound, 0x01, func() packets.Packet { return &EncryptionRequest{} })
	packets.Register(packets.Login, packets.ClientBound, 0x02, func() packets.Packet { return &LoginSuccess{} })
	packets.Register(packets.Login, packets.ClientBound, 0x03, func() packets.Packet { return &SetInitialCompression{} })
}

func (t *LoginDisconnect) Id() byte {
	return 0x00 // 0
}

func (t *LoginDisconnect) Encode(ww io.Writer) (err error) {
	// Encoding: Reason (packets.Chat)
	// Unable to encode: t.Reason (packets.Chat)

	return
}

func (t *LoginDisconnect) Decode(rr io.Reader) (err error) {
	// Decoding: Reason (packets.Chat)
	// Unable to decode: t.Reason (packets.Chat)

	return
}

func (t *EncryptionRequest) Id() byte {
	return 0x01 // 1
}

func (t *EncryptionRequest) Encode(ww io.Writer) (err error) {
	// Encoding: ServerID (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.ServerID)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: PublicKey ()
	// Unable to encode: t.PublicKey ()

	// Encoding: VerifyToken ()
	// Unable to encode: t.VerifyToken ()

	return
}

func (t *EncryptionRequest) Decode(rr io.Reader) (err error) {
	// Decoding: ServerID (string)
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
	t.ServerID = string(tmp1)

	// Decoding: PublicKey ()
	// Unable to decode: t.PublicKey ()

	// Decoding: VerifyToken ()
	// Unable to decode: t.VerifyToken ()

	return
}

func (t *LoginSuccess) Id() byte {
	return 0x02 // 2
}

func (t *LoginSuccess) Encode(ww io.Writer) (err error) {
	// Encoding: UUID (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.UUID)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: Username (string)
	tmp3 := make([]byte, binary.MaxVarintLen32)
	tmp4 := []byte(t.Username)
	tmp5 := packets.PutVarint(tmp3, int64(len(tmp4)))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp5]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp4); err != nil {
		return
	}

	return
}

func (t *LoginSuccess) Decode(rr io.Reader) (err error) {
	// Decoding: UUID (string)
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
	t.UUID = string(tmp1)

	// Decoding: Username (string)
	tmp3, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp4 := make([]byte, tmp3)
	tmp5, err := rr.Read(tmp4)
	if err != nil {
		return
	} else if int64(tmp5) != tmp3 {
		return errors.New("didn't read enough bytes for string")
	}
	t.Username = string(tmp4)

	return
}

func (t *SetInitialCompression) Id() byte {
	return 0x03 // 3
}

func (t *SetInitialCompression) Encode(ww io.Writer) (err error) {
	// Encoding: Threshold (packets.VarInt)
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(t.Threshold))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return
	}

	return
}

func (t *SetInitialCompression) Decode(rr io.Reader) (err error) {
	// Decoding: Threshold (packets.VarInt)
	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.Threshold = packets.VarInt(tmp0)

	return
}
