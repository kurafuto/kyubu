package packets

import "fmt"

type Message struct {
	PacketId byte
	PlayerId int8
	Message  string
}

func (p Message) Id() byte {
	return 0x0d
}

func (p Message) Size() int {
	return ReflectSize(p)
}

func (p Message) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadMessage(b []byte) (Packet, error) {
	var p Message
	err := ReflectRead(b, &p)
	return &p, err
}

func NewMessage(playerId int8, message string) (p *Message, err error) {
	if len(message) > StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
	}
	p = &Message{
		PacketId: 0x0d,
		PlayerId: playerId,
		Message:  message,
	}
	return
}

func init() {
	MustRegister(&PacketInfo{
		Id:   0x0d,
		Read: ReadMessage,
		Size: ReflectSize(&Message{}),
		Type: Both,
		Name: "Message",
	})
}
