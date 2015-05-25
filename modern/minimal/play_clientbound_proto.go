// ARR, HERE BE DRAGONS! DO NOT EDIT
// protocol_generator -file=play_clientbound.go -direction=clientbound -state=play -package=minimal

package minimal

import (
	"encoding/binary"
	"errors"
	"github.com/kurafuto/kyubu/packets"
	"io"
)

func init() {
	packets.Register(packets.Play, packets.ClientBound, 0x01, func() packets.Packet { return &JoinGame{} })
	packets.Register(packets.Play, packets.ClientBound, 0x02, func() packets.Packet { return &ServerMessage{} })
	packets.Register(packets.Play, packets.ClientBound, 0x07, func() packets.Packet { return &Respawn{} })
	packets.Register(packets.Play, packets.ClientBound, 0x38, func() packets.Packet { return &PlayerListItem{} })
	packets.Register(packets.Play, packets.ClientBound, 0x3a, func() packets.Packet { return &ServerTabComplete{} })
	packets.Register(packets.Play, packets.ClientBound, 0x3b, func() packets.Packet { return &ScoreboardObjective{} })
	packets.Register(packets.Play, packets.ClientBound, 0x3c, func() packets.Packet { return &UpdateScore{} })
	packets.Register(packets.Play, packets.ClientBound, 0x3d, func() packets.Packet { return &ShowScoreboard{} })
	packets.Register(packets.Play, packets.ClientBound, 0x3e, func() packets.Packet { return &Teams{} })
	packets.Register(packets.Play, packets.ClientBound, 0x3f, func() packets.Packet { return &ServerPluginMessage{} })
	packets.Register(packets.Play, packets.ClientBound, 0x40, func() packets.Packet { return &Disconnect{} })
	packets.Register(packets.Play, packets.ClientBound, 0x46, func() packets.Packet { return &SetCompression{} })
}

func (t *JoinGame) Id() byte {
	return 0x01 // 1
}

func (t *JoinGame) Encode(ww io.Writer) (err error) {
	// Encoding: EntityID (int32)
	if err = binary.Write(ww, binary.BigEndian, t.EntityID); err != nil {
		return
	}

	// Encoding: Gamemode (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.Gamemode); err != nil {
		return
	}

	// Encoding: Dimension (int8)
	if err = binary.Write(ww, binary.BigEndian, t.Dimension); err != nil {
		return
	}

	// Encoding: Difficulty (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.Difficulty); err != nil {
		return
	}

	// Encoding: MaxPlayers (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.MaxPlayers); err != nil {
		return
	}

	// Encoding: LevelType (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.LevelType)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: ReducedDebug (bool)
	tmp3 := byte(0)
	if t.ReducedDebug {
		tmp3 = byte(1)
	}
	if err = binary.Write(ww, binary.BigEndian, tmp3); err != nil {
		return
	}

	return
}

func (t *JoinGame) Decode(rr io.Reader) (err error) {
	// Decoding: EntityID (int32)
	if err = binary.Read(rr, binary.BigEndian, t.EntityID); err != nil {
		return
	}

	// Decoding: Gamemode (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.Gamemode); err != nil {
		return
	}

	// Decoding: Dimension (int8)
	if err = binary.Read(rr, binary.BigEndian, t.Dimension); err != nil {
		return
	}

	// Decoding: Difficulty (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.Difficulty); err != nil {
		return
	}

	// Decoding: MaxPlayers (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.MaxPlayers); err != nil {
		return
	}

	// Decoding: LevelType (string)
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
	t.LevelType = string(tmp1)

	// Decoding: ReducedDebug (bool)
	var tmp3 [1]byte
	if _, err = rr.Read(tmp3[:1]); err != nil {
		return
	}
	t.ReducedDebug = tmp3[0] == 0x01

	return
}

func (t *ServerMessage) Id() byte {
	return 0x02 // 2
}

func (t *ServerMessage) Encode(ww io.Writer) (err error) {
	// Encoding: Data (packets.Chat)
	// Unable to encode: t.Data (packets.Chat)

	// Encoding: Position (int8)
	if err = binary.Write(ww, binary.BigEndian, t.Position); err != nil {
		return
	}

	return
}

func (t *ServerMessage) Decode(rr io.Reader) (err error) {
	// Decoding: Data (packets.Chat)
	// Unable to decode: t.Data (packets.Chat)

	// Decoding: Position (int8)
	if err = binary.Read(rr, binary.BigEndian, t.Position); err != nil {
		return
	}

	return
}

func (t *Respawn) Id() byte {
	return 0x07 // 7
}

func (t *Respawn) Encode(ww io.Writer) (err error) {
	// Encoding: Dimension (int32)
	if err = binary.Write(ww, binary.BigEndian, t.Dimension); err != nil {
		return
	}

	// Encoding: Difficulty (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.Difficulty); err != nil {
		return
	}

	// Encoding: Gamemode (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.Gamemode); err != nil {
		return
	}

	// Encoding: LevelType (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.LevelType)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	return
}

func (t *Respawn) Decode(rr io.Reader) (err error) {
	// Decoding: Dimension (int32)
	if err = binary.Read(rr, binary.BigEndian, t.Dimension); err != nil {
		return
	}

	// Decoding: Difficulty (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.Difficulty); err != nil {
		return
	}

	// Decoding: Gamemode (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.Gamemode); err != nil {
		return
	}

	// Decoding: LevelType (string)
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
	t.LevelType = string(tmp1)

	return
}

func (t *PlayerListItem) Id() byte {
	return 0x38 // 56
}

func (t *PlayerListItem) Encode(ww io.Writer) (err error) {
	// Encoding: Action (packets.VarInt)
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(t.Action))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return
	}

	// Encoding: NumPlayers (packets.VarInt)
	tmp2 := make([]byte, binary.MaxVarintLen64)
	tmp3 := packets.PutVarint(tmp2, int64(t.NumPlayers))
	if err = binary.Write(ww, binary.BigEndian, tmp2[:tmp3]); err != nil {
		return
	}

	// Encoding: Players ()
	// Unable to encode: t.Players ()

	return
}

func (t *PlayerListItem) Decode(rr io.Reader) (err error) {
	// Decoding: Action (packets.VarInt)
	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.Action = packets.VarInt(tmp0)

	// Decoding: NumPlayers (packets.VarInt)
	tmp1, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.NumPlayers = packets.VarInt(tmp1)

	// Decoding: Players ()
	// Unable to decode: t.Players ()

	return
}

func (t *ServerTabComplete) Id() byte {
	return 0x3a // 58
}

func (t *ServerTabComplete) Encode(ww io.Writer) (err error) {
	// Encoding: Matches ()
	// Unable to encode: t.Matches ()

	return
}

func (t *ServerTabComplete) Decode(rr io.Reader) (err error) {
	// Decoding: Matches ()
	// Unable to decode: t.Matches ()

	return
}

func (t *ScoreboardObjective) Id() byte {
	return 0x3b // 59
}

func (t *ScoreboardObjective) Encode(ww io.Writer) (err error) {
	// Encoding: Name (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Name)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: Mode (uint8)
	if err = binary.Write(ww, binary.BigEndian, t.Mode); err != nil {
		return
	}

	// Encoding: Value (string)
	tmp3 := make([]byte, binary.MaxVarintLen32)
	tmp4 := []byte(t.Value)
	tmp5 := packets.PutVarint(tmp3, int64(len(tmp4)))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp5]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp4); err != nil {
		return
	}

	// Encoding: Type (string)
	tmp6 := make([]byte, binary.MaxVarintLen32)
	tmp7 := []byte(t.Type)
	tmp8 := packets.PutVarint(tmp6, int64(len(tmp7)))
	if err = binary.Write(ww, binary.BigEndian, tmp6[:tmp8]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp7); err != nil {
		return
	}

	return
}

func (t *ScoreboardObjective) Decode(rr io.Reader) (err error) {
	// Decoding: Name (string)
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
	t.Name = string(tmp1)

	// Decoding: Mode (uint8)
	if err = binary.Read(rr, binary.BigEndian, t.Mode); err != nil {
		return
	}

	// Decoding: Value (string)
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
	t.Value = string(tmp4)

	// Decoding: Type (string)
	tmp6, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp7 := make([]byte, tmp6)
	tmp8, err := rr.Read(tmp7)
	if err != nil {
		return
	} else if int64(tmp8) != tmp6 {
		return errors.New("didn't read enough bytes for string")
	}
	t.Type = string(tmp7)

	return
}

func (t *UpdateScore) Id() byte {
	return 0x3c // 60
}

func (t *UpdateScore) Encode(ww io.Writer) (err error) {
	// Encoding: Name (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Name)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: Action (int8)
	if err = binary.Write(ww, binary.BigEndian, t.Action); err != nil {
		return
	}

	// Encoding: ObjectiveName (string)
	tmp3 := make([]byte, binary.MaxVarintLen32)
	tmp4 := []byte(t.ObjectiveName)
	tmp5 := packets.PutVarint(tmp3, int64(len(tmp4)))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp5]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp4); err != nil {
		return
	}

	// Encoding: Value (packets.VarInt)
	tmp6 := make([]byte, binary.MaxVarintLen64)
	tmp7 := packets.PutVarint(tmp6, int64(t.Value))
	if err = binary.Write(ww, binary.BigEndian, tmp6[:tmp7]); err != nil {
		return
	}

	return
}

func (t *UpdateScore) Decode(rr io.Reader) (err error) {
	// Decoding: Name (string)
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
	t.Name = string(tmp1)

	// Decoding: Action (int8)
	if err = binary.Read(rr, binary.BigEndian, t.Action); err != nil {
		return
	}

	// Decoding: ObjectiveName (string)
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
	t.ObjectiveName = string(tmp4)

	// Decoding: Value (packets.VarInt)
	tmp6, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.Value = packets.VarInt(tmp6)

	return
}

func (t *ShowScoreboard) Id() byte {
	return 0x3d // 61
}

func (t *ShowScoreboard) Encode(ww io.Writer) (err error) {
	// Encoding: Position (int8)
	if err = binary.Write(ww, binary.BigEndian, t.Position); err != nil {
		return
	}

	// Encoding: Name (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Name)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	return
}

func (t *ShowScoreboard) Decode(rr io.Reader) (err error) {
	// Decoding: Position (int8)
	if err = binary.Read(rr, binary.BigEndian, t.Position); err != nil {
		return
	}

	// Decoding: Name (string)
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
	t.Name = string(tmp1)

	return
}

func (t *Teams) Id() byte {
	return 0x3e // 62
}

func (t *Teams) Encode(ww io.Writer) (err error) {
	// Encoding: Name (string)
	tmp0 := make([]byte, binary.MaxVarintLen32)
	tmp1 := []byte(t.Name)
	tmp2 := packets.PutVarint(tmp0, int64(len(tmp1)))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp2]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp1); err != nil {
		return
	}

	// Encoding: Mode (int8)
	if err = binary.Write(ww, binary.BigEndian, t.Mode); err != nil {
		return
	}

	// Encoding: Display (string)
	tmp3 := make([]byte, binary.MaxVarintLen32)
	tmp4 := []byte(t.Display)
	tmp5 := packets.PutVarint(tmp3, int64(len(tmp4)))
	if err = binary.Write(ww, binary.BigEndian, tmp3[:tmp5]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp4); err != nil {
		return
	}

	// Encoding: Prefix (string)
	tmp6 := make([]byte, binary.MaxVarintLen32)
	tmp7 := []byte(t.Prefix)
	tmp8 := packets.PutVarint(tmp6, int64(len(tmp7)))
	if err = binary.Write(ww, binary.BigEndian, tmp6[:tmp8]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp7); err != nil {
		return
	}

	// Encoding: Suffix (string)
	tmp9 := make([]byte, binary.MaxVarintLen32)
	tmp10 := []byte(t.Suffix)
	tmp11 := packets.PutVarint(tmp9, int64(len(tmp10)))
	if err = binary.Write(ww, binary.BigEndian, tmp9[:tmp11]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp10); err != nil {
		return
	}

	// Encoding: FriendlyFire (int8)
	if err = binary.Write(ww, binary.BigEndian, t.FriendlyFire); err != nil {
		return
	}

	// Encoding: NameTagVisibility (string)
	tmp12 := make([]byte, binary.MaxVarintLen32)
	tmp13 := []byte(t.NameTagVisibility)
	tmp14 := packets.PutVarint(tmp12, int64(len(tmp13)))
	if err = binary.Write(ww, binary.BigEndian, tmp12[:tmp14]); err != nil {
		return
	}
	if err = binary.Write(ww, binary.BigEndian, tmp13); err != nil {
		return
	}

	// Encoding: Color (int8)
	if err = binary.Write(ww, binary.BigEndian, t.Color); err != nil {
		return
	}

	// Encoding: Players ()
	// Unable to encode: t.Players ()

	return
}

func (t *Teams) Decode(rr io.Reader) (err error) {
	// Decoding: Name (string)
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
	t.Name = string(tmp1)

	// Decoding: Mode (int8)
	if err = binary.Read(rr, binary.BigEndian, t.Mode); err != nil {
		return
	}

	// Decoding: Display (string)
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
	t.Display = string(tmp4)

	// Decoding: Prefix (string)
	tmp6, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp7 := make([]byte, tmp6)
	tmp8, err := rr.Read(tmp7)
	if err != nil {
		return
	} else if int64(tmp8) != tmp6 {
		return errors.New("didn't read enough bytes for string")
	}
	t.Prefix = string(tmp7)

	// Decoding: Suffix (string)
	tmp9, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp10 := make([]byte, tmp9)
	tmp11, err := rr.Read(tmp10)
	if err != nil {
		return
	} else if int64(tmp11) != tmp9 {
		return errors.New("didn't read enough bytes for string")
	}
	t.Suffix = string(tmp10)

	// Decoding: FriendlyFire (int8)
	if err = binary.Read(rr, binary.BigEndian, t.FriendlyFire); err != nil {
		return
	}

	// Decoding: NameTagVisibility (string)
	tmp12, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	tmp13 := make([]byte, tmp12)
	tmp14, err := rr.Read(tmp13)
	if err != nil {
		return
	} else if int64(tmp14) != tmp12 {
		return errors.New("didn't read enough bytes for string")
	}
	t.NameTagVisibility = string(tmp13)

	// Decoding: Color (int8)
	if err = binary.Read(rr, binary.BigEndian, t.Color); err != nil {
		return
	}

	// Decoding: Players ()
	// Unable to decode: t.Players ()

	return
}

func (t *ServerPluginMessage) Id() byte {
	return 0x3f // 63
}

func (t *ServerPluginMessage) Encode(ww io.Writer) (err error) {
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

func (t *ServerPluginMessage) Decode(rr io.Reader) (err error) {
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

func (t *Disconnect) Id() byte {
	return 0x40 // 64
}

func (t *Disconnect) Encode(ww io.Writer) (err error) {
	// Encoding: Reason (packets.Chat)
	// Unable to encode: t.Reason (packets.Chat)

	return
}

func (t *Disconnect) Decode(rr io.Reader) (err error) {
	// Decoding: Reason (packets.Chat)
	// Unable to decode: t.Reason (packets.Chat)

	return
}

func (t *SetCompression) Id() byte {
	return 0x46 // 70
}

func (t *SetCompression) Encode(ww io.Writer) (err error) {
	// Encoding: Threshold (packets.VarInt)
	tmp0 := make([]byte, binary.MaxVarintLen64)
	tmp1 := packets.PutVarint(tmp0, int64(t.Threshold))
	if err = binary.Write(ww, binary.BigEndian, tmp0[:tmp1]); err != nil {
		return
	}

	return
}

func (t *SetCompression) Decode(rr io.Reader) (err error) {
	// Decoding: Threshold (packets.VarInt)
	tmp0, err := packets.ReadVarint(rr)
	if err != nil {
		return
	}
	t.Threshold = packets.VarInt(tmp0)

	return
}
