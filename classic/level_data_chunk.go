package classic

import "github.com/sysr-q/kyubu/packets"

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
	return packets.ReflectSize(p)
}

func (p LevelDataChunk) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadLevelDataChunk(b []byte) (packets.Packet, error) {
	var p LevelDataChunk
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewLevelDataChunk(chunkData []byte, percentComplete byte) (p *LevelDataChunk, err error) {
	if len(chunkData) > packets.BytesSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes", packets.BytesSize)
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
	packets.Register(&packets.PacketInfo{
		Id:        0x03,
		Read:      ReadLevelDataChunk,
		Size:      packets.ReflectSize(&LevelDataChunk{}),
		Direction: packets.Anomalous,
		Name:      "Level Data Chunk",
	})
}
