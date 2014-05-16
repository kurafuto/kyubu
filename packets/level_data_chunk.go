package packets

import "fmt"

type LevelDataChunk struct {
	PacketId        byte
	ChunkLength     int16
	ChunkData       []byte
	PercentComplete byte
}

func (p LevelDataChunk) Id() byte {
	return 0x03
}

func (p LevelDataChunk) Size() int {
	return ReflectSize(p)
}

func (p LevelDataChunk) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadLevelDataChunk(b []byte) (Packet, error) {
	var p LevelDataChunk
	err := ReflectRead(b, &p)
	return &p, err
}

func NewLevelDataChunk(chunkData []byte, percentComplete byte) (p *LevelDataChunk, err error) {
	if len(chunkData) > BytesSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes", BytesSize)
	}
	p = &LevelDataChunk{
		PacketId:        0x03,
		ChunkLength:     int16(len(chunkData)),
		ChunkData:       chunkData,
		PercentComplete: percentComplete,
	}
	return
}

func init() {
	MustRegister(&PacketInfo{
		Id:   0x03,
		Read: ReadLevelDataChunk,
		Size: ReflectSize(&LevelDataChunk{}),
		Type: ServerOnly,
		Name: "Level Data Chunk",
	})
}
