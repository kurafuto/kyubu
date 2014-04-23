package packets

const MessageSize = (ByteSize + // Packet ID
	SByteSize + // Player ID
	StringSize) // Message

type Message struct {
	PacketId byte
	PlayerId int8
	Message  string
}

func (p *Message) Id() byte {
	return p.PacketId
}

func (p *Message) Size() int {
	return MessageSize
}

func (p *Message) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteSByte(p.PlayerId)
	raw.WriteString(p.Message)
	return raw.Buffer.Bytes()
}

func ReadMessage(b []byte) (Packet, error) {
	p := Message{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if playerId, err := raw.ReadSByte(); err != nil {
		return nil, err
	} else {
		p.PlayerId = playerId
	}
	if message, err := raw.ReadString(); err != nil {
		return nil, err
	} else {
		p.Message = message
	}
	return &p, nil
}

func NewMessage(playerId int8, message string) (p *Message, err error) {
	p = &Message{
		PacketId: 13,
		PlayerId: playerId,
		Message:  message,
	}
	return
}
