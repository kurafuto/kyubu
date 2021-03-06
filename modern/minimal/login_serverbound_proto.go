// Generated by protocol_generator; DO NOT EDIT
// protocol_generator -file=login_serverbound.go -direction=serverbound -state=login -package=modern

package modern

import (
	"encoding/binary"
	"fmt"
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
	if err = packets.WriteString(ww, t.Username); err != nil {
		return err
	}
	return
}

func (t *LoginStart) Decode(rr io.Reader) (err error) {
	if t.Username, err = packets.ReadString(rr); err != nil {
		return err
	}
	return
}

func (t *EncryptionResponse) Id() byte {
	return 0x01 // 1
}

func (t *EncryptionResponse) Encode(ww io.Writer) (err error) {
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(len(t.SharedSecret)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return err
	}

	if _, err = ww.Write(t.SharedSecret); err != nil {
		return err
	}
	tmp2 := make([]byte, binary.MaxVarintLen64)
	tmp3 := packets.PutVarint(tmp2, int64(len(t.VerifyToken)))
	if err = binary.Write(ww, binary.BigEndian, tmp2[:tmp3]); err != nil {
		return err
	}

	if _, err = ww.Write(t.VerifyToken); err != nil {
		return err
	}
	return
}

func (t *EncryptionResponse) Decode(rr io.Reader) (err error) {
	var tmp0 packets.VarInt
	tmp1, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	tmp0 = packets.VarInt(tmp1)

	if tmp0 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp0)
	}
	t.SharedSecret = make([]byte, tmp0)
	if _, err = rr.Read(t.SharedSecret); err != nil {
		return err
	}
	var tmp2 packets.VarInt
	tmp3, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	tmp2 = packets.VarInt(tmp3)

	if tmp2 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp2)
	}
	t.VerifyToken = make([]byte, tmp2)
	if _, err = rr.Read(t.VerifyToken); err != nil {
		return err
	}
	return
}
