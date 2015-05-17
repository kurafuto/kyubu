package classic

import "github.com/sysr-q/kyubu/packets"

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
	return packets.ReflectSize(p)
}

func (p Message) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadMessage(b []byte) (packets.Packet, error) {
	var p Message
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewMessage(playerId int8, message string) (p *Message, err error) {
	if len(message) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &Message{
		PacketId: 0x0d,
		PlayerId: playerId,
		Message:  message,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x0d,
		Read:      ReadMessage,
		Size:      packets.ReflectSize(&Message{}),
		Direction: packets.Anomalous,
		Name:      "Message",
	})
}
