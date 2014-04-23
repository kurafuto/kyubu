package packets

import "fmt"

const LevelDataChunkSize = (ByteSize + // Packet ID
	ShortSize + // Chunk Length
	BytesSize + // Chunk Data
	ByteSize) // Percent Complete

type LevelDataChunk struct {
	PacketId        byte
	ChunkLength     int16
	ChunkData       []byte
	PercentComplete byte
	// Not packed, flag indicating whether implementations have
	// internally handled this packet before pushing it through.
	Handled bool
}

func (p *LevelDataChunk) Id() byte {
	return p.PacketId
}

func (p *LevelDataChunk) Size() int {
	return LevelDataChunkSize
}

func (p *LevelDataChunk) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteShort(p.ChunkLength)
	raw.WriteBytes(p.ChunkData)
	raw.WriteByte(p.PercentComplete)
	return raw.Buffer.Bytes()
}

func ReadLevelDataChunk(b []byte) (Packet, error) {
	p := LevelDataChunk{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if chunkLength, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.ChunkLength = chunkLength
	}
	if chunkData, err := raw.ReadBytes(); err != nil {
		return nil, err
	} else {
		// TODO: decompress gzip
		p.ChunkData = chunkData
	}
	if percentComplete, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PercentComplete = percentComplete
	}
	return &p, nil
}

func NewLevelDataChunk(chunkData []byte, percentComplete byte) (p *LevelDataChunk, err error) {
	if len(chunkData) > BytesSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes", BytesSize)
	}
	p = &LevelDataChunk{
		PacketId:        3,
		ChunkLength:     int16(len(chunkData)),
		ChunkData:       chunkData,
		PercentComplete: percentComplete,
		Handled:         false,
	}
	return
}
