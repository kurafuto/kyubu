package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

var Endianness = binary.BigEndian

// Size constants from spec, used to deduce packet size.
const (
	ByteSize   = 1
	SByteSize  = 1
	ShortSize  = 2
	StringSize = 64
	BytesSize  = 1024
)

const (
	ServerOnly = iota
	ClientOnly
	Both
)

type Packet interface {
	Id() byte
	Size() int
	Bytes() []byte
}

type PacketInfo struct {
	Read  func([]byte) (Packet, error)
	Size  int
	Type  int
	Ident string
}

var Packets = map[byte]PacketInfo{
	0x00: PacketInfo{
		Read:  ReadIdentification,
		Size:  IdentificationSize,
		Type:  Both,
		Ident: "Identification",
	},
	0x01: PacketInfo{
		Read:  ReadPing,
		Size:  PingSize,
		Type:  ServerOnly,
		Ident: "Ping",
	},
	0x02: PacketInfo{
		Read:  ReadLevelInitialize,
		Size:  LevelInitializeSize,
		Type:  ServerOnly,
		Ident: "Level Initialize",
	},
	0x03: PacketInfo{
		Read:  ReadLevelDataChunk,
		Size:  LevelDataChunkSize,
		Type:  ServerOnly,
		Ident: "Level Data Chunk",
	},
	0x04: PacketInfo{
		Read:  ReadLevelFinalize,
		Size:  LevelFinalizeSize,
		Type:  ServerOnly,
		Ident: "Level Finalize",
	},
	0x05: PacketInfo{
		Read:  ReadSetBlock5,
		Size:  SetBlock5Size,
		Type:  ClientOnly,
		Ident: "Set Block [Client]",
	},
	0x06: PacketInfo{
		Read:  ReadSetBlock6,
		Size:  SetBlock6Size,
		Type:  ServerOnly,
		Ident: "Set Block [Server]",
	},
	0x07: PacketInfo{
		Read:  ReadSpawnPlayer,
		Size:  SpawnPlayerSize,
		Type:  ServerOnly,
		Ident: "Spawn Player",
	},
	0x08: PacketInfo{
		Read:  ReadPositionOrientation,
		Size:  PositionOrientationSize,
		Type:  Both,
		Ident: "Position & Orientation",
	},
	0x09: PacketInfo{
		Read:  ReadPositionOrientationUpdate,
		Size:  PositionOrientationUpdateSize,
		Type:  ServerOnly,
		Ident: "Position & Orientation Update",
	},
	0x0a: PacketInfo{
		Read:  ReadPositionUpdate,
		Size:  PositionUpdateSize,
		Type:  ServerOnly,
		Ident: "Position Update",
	},
	0x0b: PacketInfo{
		Read:  ReadOrientationUpdate,
		Size:  OrientationUpdateSize,
		Type:  ServerOnly,
		Ident: "Orientation Update",
	},
	0x0c: PacketInfo{
		Read:  ReadDespawnPlayer,
		Size:  DespawnPlayerSize,
		Type:  ServerOnly,
		Ident: "Despawn Player",
	},
	0x0d: PacketInfo{
		Read:  ReadMessage,
		Size:  MessageSize,
		Type:  Both,
		Ident: "Message",
	},
	0x0e: PacketInfo{
		Read:  ReadDisconnectPlayer,
		Size:  DisconnectPlayerSize,
		Type:  ServerOnly,
		Ident: "Disconnect Player",
	},
	0x0f: PacketInfo{
		Read:  ReadUpdateUserType,
		Size:  UpdateUserTypeSize,
		Type:  ServerOnly,
		Ident: "Update User Type",
	},
}

// PacketWrapper wraps an array of bytes into a bytes.Buffer, and exposes
// various Read/Write type helpers. It also optionally (and by default) strips
// padding characters, as defined by the Minecraft Classic Protocol spec.
type PacketWrapper struct {
	Buffer *bytes.Buffer
	Strip  bool
}

// Write functions - append data to the buffer.

func (p *PacketWrapper) WriteByte(b byte) error {
	return binary.Write(p.Buffer, Endianness, &b)
}

func (p *PacketWrapper) WriteSByte(i int8) error {
	return binary.Write(p.Buffer, Endianness, &i)
}

func (p *PacketWrapper) WriteShort(i int16) error {
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

// Read functions - read data out of the buffer.

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
