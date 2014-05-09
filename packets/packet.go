package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

var (
	Endianness    = binary.BigEndian
	AllowOverride = false // Can we replace packets?
)

// Size constants from spec, used to deduce packet size.
const (
	ByteSize   = 1
	SByteSize  = 1
	ShortSize  = 2
	IntSize    = 4
	StringSize = 64
	BytesSize  = 1024
)

// PacketType lets us define whether packets are purely C->S, S->C or C<->S.
// It's not an important implementation detail, but it's nice to know.
type PacketType byte

func (t PacketType) String() string {
	if t == ServerOnly {
		return "S->C"
	} else if t == ClientOnly {
		return "C->S"
	} else if t == Both {
		return "C<>S"
	}
	return "????"
}

const (
	ServerOnly PacketType = iota
	ClientOnly
	Both
)

type Packet interface {
	Id() byte
	Size() int
	Bytes() []byte
}

type PacketInfo struct {
	Id   byte
	Read func([]byte) (Packet, error)
	Size int
	Type PacketType
	Name string // Human presentable
}

type readFunc func([]byte) (Packet, error)

var Packets = map[byte]*PacketInfo{}

// Register allows users to register additional packets with the internal parser.
// These packets will be recognised, parsed and sent to the parser channel.
func Register(p *PacketInfo) (bool, error) {
	if _, ok := Packets[p.Id]; ok && !AllowOverride {
		return false, fmt.Errorf("Packet already registered to id %#.2x", p.Id)
	}
	Packets[p.Id] = p
	return true, nil
}

// MustRegister is functionally the same as Register, however if an error is
// returned, or the registration fails for whatever reason, it will panic.
func MustRegister(p *PacketInfo) {
	ok, err := Register(p)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic(fmt.Errorf("kyubu/packets: Registration of new packet %#.2x failed", p.Id))
	}
}

func init() {
	Register(&PacketInfo{
		Id:   0x00,
		Read: ReadIdentification,
		Size: IdentificationSize,
		Type: Both,
		Name: "Identification",
	})
	Register(&PacketInfo{
		Id:   0x01,
		Read: ReadPing,
		Size: PingSize,
		Type: ServerOnly,
		Name: "Ping",
	})
	Register(&PacketInfo{
		Id:   0x02,
		Read: ReadLevelInitialize,
		Size: LevelInitializeSize,
		Type: ServerOnly,
		Name: "Level Initialize",
	})
	Register(&PacketInfo{
		Id:   0x03,
		Read: ReadLevelDataChunk,
		Size: LevelDataChunkSize,
		Type: ServerOnly,
		Name: "Level Data Chunk",
	})
	Register(&PacketInfo{
		Id:   0x04,
		Read: ReadLevelFinalize,
		Size: LevelFinalizeSize,
		Type: ServerOnly,
		Name: "Level Finalize",
	})
	Register(&PacketInfo{
		Id:   0x05,
		Read: ReadSetBlock5,
		Size: SetBlock5Size,
		Type: ClientOnly,
		Name: "Set Block [C->S]",
	})
	Register(&PacketInfo{
		Id:   0x06,
		Read: ReadSetBlock6,
		Size: SetBlock6Size,
		Type: ServerOnly,
		Name: "Set Block [S->C]",
	})
	Register(&PacketInfo{
		Id:   0x07,
		Read: ReadSpawnPlayer,
		Size: SpawnPlayerSize,
		Type: ServerOnly,
		Name: "Spawn Player",
	})
	Register(&PacketInfo{
		Id:   0x08,
		Read: ReadPositionOrientation,
		Size: PositionOrientationSize,
		Type: Both,
		Name: "Position/Orientation",
	})
	Register(&PacketInfo{
		Id:   0x09,
		Read: ReadPositionOrientationUpdate,
		Size: PositionOrientationUpdateSize,
		Type: ServerOnly,
		Name: "Position/Orientation Update",
	})
	Register(&PacketInfo{
		Id:   0x0a,
		Read: ReadPositionUpdate,
		Size: PositionUpdateSize,
		Type: ServerOnly,
		Name: "Position Update",
	})
	Register(&PacketInfo{
		Id:   0x0b,
		Read: ReadOrientationUpdate,
		Size: OrientationUpdateSize,
		Type: ServerOnly,
		Name: "Orientation Update",
	})
	Register(&PacketInfo{
		Id:   0x0c,
		Read: ReadDespawnPlayer,
		Size: DespawnPlayerSize,
		Type: ServerOnly,
		Name: "Despawn Player",
	})
	Register(&PacketInfo{
		Id:   0x0d,
		Read: ReadMessage,
		Size: MessageSize,
		Type: Both,
		Name: "Message",
	})
	Register(&PacketInfo{
		Id:   0x0e,
		Read: ReadDisconnectPlayer,
		Size: DisconnectPlayerSize,
		Type: ServerOnly,
		Name: "Disconnect Player",
	})
	Register(&PacketInfo{
		Id:   0x0f,
		Read: ReadUpdateUserType,
		Size: UpdateUserTypeSize,
		Type: ServerOnly,
		Name: "Update User Type",
	})
}

// PacketWrapper wraps an array of bytes into a bytes.Buffer, and exposes
// various Read/Write type helpers. It also optionally (and by default) strips
// padding characters, as defined by the Minecraft Classic Protocol spec.
type PacketWrapper struct {
	Buffer *bytes.Buffer
	Strip  bool
}

func (p *PacketWrapper) WriteByte(b byte) error {
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) WriteSByte(i int8) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteShort(i int16) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteInt(i int32) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteString(s string) error {
	if len(s) > StringSize {
		return fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
	}
	// Pad with spaces, up to 64 bytes.
	s += strings.Repeat(" ", StringSize-len(s))
	b := []byte(s)
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) WriteBytes(b []byte) error {
	if len(b) > BytesSize {
		return fmt.Errorf("kyubu/packets: cannot write over %d bytes", BytesSize)
	}
	// Pad with 0x00, up to 1024 bytes.
	b = append(b, bytes.Repeat([]byte{0}, BytesSize-len(b))...)
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) ReadByte() (i byte, err error) {
	b := make([]byte, 1)
	_, err = p.Buffer.Read(b)
	if err != nil {
		return
	}
	i = b[0]
	return
}

func (p *PacketWrapper) ReadSByte() (i int8, err error) {
	err = binary.Read(p.Buffer, Endianness, &i)
	return
}

func (p *PacketWrapper) ReadShort() (i int16, err error) {
	err = binary.Read(p.Buffer, Endianness, &i)
	return
}

func (p *PacketWrapper) ReadInt() (i int32, err error) {
	err = binary.Read(p.Buffer, Endianness, &i)
	return
}

func (p *PacketWrapper) ReadString() (s string, err error) {
	// Strings are 64 bytes, padded with spaces (0x20).
	b := make([]byte, StringSize)
	_, err = p.Buffer.Read(b)
	if err != nil {
		return
	}
	s = string(b)
	if p.Strip {
		s = strings.TrimRight(s, " ")
	}
	return
}

func (p *PacketWrapper) ReadBytes() (b []byte, err error) {
	b = make([]byte, BytesSize)
	_, err = p.Buffer.Read(b)
	if p.Strip {
		b = bytes.TrimRight(b, string([]byte{0}))
	}
	return
}

func NewPacketWrapper(b []byte) (p *PacketWrapper) {
	p = &PacketWrapper{
		Buffer: bytes.NewBuffer(b),
		Strip:  true,
	}
	return
}
