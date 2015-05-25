// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=play_serverbound.go -direction=serverbound -state=play -package=minimal

package minimal

import (
	"encoding/binary"
	"errors"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Play, packets.ServerBound, 0x01, func() packets.Packet { return &ClientMessage{} })
	packets.Register(packets.Play, packets.ServerBound, 0x14, func() packets.Packet { return &ClientTabComplete{} })
	packets.Register(packets.Play, packets.ServerBound, 0x16, func() packets.Packet { return &ClientStatus{} })
	packets.Register(packets.Play, packets.ServerBound, 0x17, func() packets.Packet { return &ClientPluginMessage{} })
}

func (t *ClientMessage) Id() byte {
	return 0x01 // 1
}

func (t *ClientMessage) Encode(ww io.Writer) (err error) {
	// Encoding: Message (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Message)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	return
}

func (t *ClientMessage) Decode(rr io.Reader) (err error) {
	// Decoding: Message (string)
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
	t.Message = string(tmp1)

	return
}

func (t *ClientTabComplete) Id() byte {
	return 0x14 // 20
}

func (t *ClientTabComplete) Encode(ww io.Writer) (err error) {
	// Encoding: Text (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Text)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: HasPosition (bool)
	tmp3 := byte(0)
	if t.HasPosition {
		tmp3 = byte(1)
	}
	if err = binary.Write(ww, binary.BigEndian, tmp3); err != nil {
		return
	}

	// Encoding: LookedAtBlock (packets.Position)
	if err = binary.Write(ww, binary.BigEndian, t.LookedAtBlock); err != nil {
		return
	}

	return
}

func (t *ClientTabComplete) Decode(rr io.Reader) (err error) {
	// Decoding: Text (string)
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
	t.Text = string(tmp1)

	// Decoding: HasPosition (bool)
	var tmp3 [1]byte
	if _, err = rr.Read(tmp3[:1]); err != nil {
		return
	}
	t.HasPosition = tmp3[0] == 0x01

	// Decoding: LookedAtBlock (packets.Position)
	if err = binary.Read(rr, binary.BigEndian, t.LookedAtBlock); err != nil {
		return
	}

	return
}

func (t *ClientStatus) Id() byte {
	return 0x16 // 22
}

func (t *ClientStatus) Encode(ww io.Writer) (err error) {
	// Encoding: ActionID (packets.VarInt)
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(t.ActionID))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return
	}

	return
}

func (t *ClientStatus) Decode(rr io.Reader) (err error) {
	// Decoding: ActionID (packets.VarInt)
	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.ActionID = packets.VarInt(tmp0)

	return
}

func (t *ClientPluginMessage) Id() byte {
	return 0x17 // 23
}

func (t *ClientPluginMessage) Encode(ww io.Writer) (err error) {
	// Encoding: Channel (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Channel)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: Data ()
	// Unable to encode: t.Data ()

	return
}

func (t *ClientPluginMessage) Decode(rr io.Reader) (err error) {
	// Decoding: Channel (string)
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
	t.Channel = string(tmp1)

	// Decoding: Data ()
	// Unable to decode: t.Data ()

	return
}
