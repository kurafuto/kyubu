// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=test.go -direction=anomalous -state=handshake -package=main

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Handshake, packets.Anomalous, 0xfe, func() packets.Packet { return &TestPacket{} })
	packets.Register(packets.Handshake, packets.Anomalous, 0xfd, func() packets.Packet { return &Test{} })
}

func (t *TestPacket) Id() byte {
	return 0xfe // 254
}

func (t *TestPacket) Encode(ww io.Writer) (err error) {
	tmp0 := byte(0)
	if t.b {
		tmp0 = byte(1)
	}
	if err = binary.Write(ww, binary.BigEndian, tmp0); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.i8); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.ui8); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.i16); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.ui16); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.i32); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.i64); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.f32); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.f64); err != nil {
		return err
	}

	tmp1 := make([]byte, binary.MaxVarintLen64)
	tmp2 := packets.PutVarint(tmp1, int64(t.V32))
	if err = binary.Write(ww, binary.BigEndian, tmp1[:tmp2]); err != nil {
		return err
	}

	tmp3 := make([]byte, binary.MaxVarintLen64)
	tmp4 := packets.PutVarint(tmp3, int64(t.V64))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp4]); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.P); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.A); err != nil {
		return err
	}

	if err = binary.Write(ww, binary.BigEndian, t.U[:]); err != nil {
		return err
	}

	return
}

func (t *TestPacket) Decode(rr io.Reader) (err error) {
	var tmp0 [1]byte
	if _, err = rr.Read(tmp0[:1]); err != nil {
		return err
	}
	t.b = tmp0[0] == 0x01

	if err = binary.Read(rr, binary.BigEndian, t.i8); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.ui8); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.i16); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.ui16); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.i32); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.i64); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.f32); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.f64); err != nil {
		return err
	}

	tmp1, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	t.V32 = packets.VarInt(tmp1)

	tmp2, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	t.V64 = packets.VarLong(tmp2)

	if err = binary.Read(rr, binary.BigEndian, t.P); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.A); err != nil {
		return err
	}

	if err = binary.Read(rr, binary.BigEndian, t.U); err != nil {
		return err
	}

	return
}

func (t *Test) Id() byte {
	return 0xfd // 253
}

func (t *Test) Encode(ww io.Writer) (err error) {
	if err = binary.Write(ww, binary.BigEndian, t.T.X); err != nil {
		return err
	}

	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.T.Y)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return err
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return err
	}

	tmp3 := make([]byte, binary.MaxVarintLen64)
	tmp4 := packets.PutVarint(tmp3, int64(t.T.Z))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp4]); err != nil {
		return err
	}

	tmp5 := make([]byte, binary.MaxVarintLen64)
	tmp6 := packets.PutVarint(tmp5, int64(len(t.TXs)))
	if err = binary.Write(ww, binary.BigEndian, tmp5[:tmp6]); err != nil {
		return err
	}

	for tmp7 := range t.TXs {
		if err = binary.Write(ww, binary.BigEndian, t.TXs[tmp7].X); err != nil {
			return err
		}

		tmp8 := make([]byte, binary.MaxVarintLen32)
		tmp9 := []byte(t.TXs[tmp7].Y)
		tmp10 := packets.PutVarint(tmp8, int64(len(tmp9)))
		if err = binary.Write(ww, binary.BigEndian, tmp8[:tmp10]); err != nil {
			return err
		}
		if err = binary.Write(ww, binary.BigEndian, tmp9); err != nil {
			return err
		}

		tmp11 := make([]byte, binary.MaxVarintLen64)
		tmp12 := packets.PutVarint(tmp11, int64(t.TXs[tmp7].Z))
		if err = binary.Write(ww, binary.BigEndian, tmp11[:tmp12]); err != nil {
			return err
		}

	}
	return
}

func (t *Test) Decode(rr io.Reader) (err error) {
	if err = binary.Read(rr, binary.BigEndian, t.T.X); err != nil {
		return err
	}

	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	tmp1 := make([]byte, tmp0)
	tmp2, err := rr.Read(tmp1)
	if err != nil {
		return err
	} else if int64(tmp2) != tmp0 {
		return errors.New("didn't read enough bytes for string")
	}
	t.T.Y = string(tmp1)

	tmp3, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	t.T.Z = packets.VarInt(tmp3)

	var tmp4 packets.VarInt
	tmp5, err := packets.ReadVarint(rr)
	if err != nil {
		return err
	}
	tmp4 = packets.VarInt(tmp5)

	if tmp4 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp4)
	}
	t.TXs = make([]TX, tmp4)
	for tmp6 := range t.TXs {
		if err = binary.Read(rr, binary.BigEndian, t.TXs[tmp6].X); err != nil {
			return err
		}

		tmp7, err := packets.ReadVarint(rr)
		if err != nil {
			return err
		}
		tmp8 := make([]byte, tmp7)
		tmp9, err := rr.Read(tmp8)
		if err != nil {
			return err
		} else if int64(tmp9) != tmp7 {
			return errors.New("didn't read enough bytes for string")
		}
		t.TXs[tmp6].Y = string(tmp8)

		tmp10, err := packets.ReadVarint(rr)
		if err != nil {
			return err
		}
		t.TXs[tmp6].Z = packets.VarInt(tmp10)

	}
	return
}
