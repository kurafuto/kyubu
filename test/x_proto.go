// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=x.go -direction=anomalous -state=handshake -package=main

package main

import (
	"encoding/binary"
	"fmt"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Handshake, packets.Anomalous, 0x42, func() packets.Packet { return &X{} })
}

func (t *X) Id() byte {
	return 0x42 // 66
}

func (t *X) Encode(ww io.Writer) (err error) {
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(len(t.Ys)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return err
	}

	for tmp2 := range t.Ys {
		tmp3 := byte(0)
		if t.Ys[tmp2].b {
			tmp3 = byte(1)
		}
		if err = binary.Write(ww, binary.BigEndian, tmp3); err != nil {
			return err
		}

	}
	return
}

func (t *X) Decode(rr io.Reader) (err error) {
	var tmp0 packets.VarInt
	tmp1, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	tmp0 = packets.VarInt(tmp1)

	if tmp0 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp0)
	}
	t.Ys = make([]Y, tmp0)
	for tmp2 := range t.Ys {
		var tmp3 [1]byte
		if _, err = rr.Read(tmp3[:1]); err != nil {
			return err
		}
		t.Ys[tmp2].b = tmp3[0] == 0x01

	}
	return
}
