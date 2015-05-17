package classic

import "github.com/sysr-q/kyubu/packets"

import "fmt"

type Identification struct {
	PacketId        byte
	ProtocolVersion byte
	Name            string
	KeyMotd         string
	UserType        byte
}

func (p Identification) Id() byte {
	return 0x00
}

func (p Identification) Size() int {
	return packets.ReflectSize(p)
}

func (p Identification) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadIdentification(b []byte) (packets.Packet, error) {
	var p Identification
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewIdentification(name, key string) (p *Identification, err error) {
	if len(name) > packets.StringSize || len(key) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", packets.StringSize)
	}
	p = &Identification{
		PacketId:        0x00,
		ProtocolVersion: 0x07,
		Name:            name,
		KeyMotd:         key,
		UserType:        0x00,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x00,
		Read:      ReadIdentification,
		Size:      packets.ReflectSize(&Identification{}),
		Direction: packets.Anomalous,
		Name:      "Identification",
	})
}
