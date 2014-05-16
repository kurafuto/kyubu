package packets

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
	return ReflectSize(p)
}

func (p Identification) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadIdentification(b []byte) (Packet, error) {
	var p Identification
	err := ReflectRead(b, &p)
	return &p, err
}

func NewIdentification(name, key string) (p *Identification, err error) {
	if len(name) > StringSize || len(key) > StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
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
	MustRegister(&PacketInfo{
		Id:   0x00,
		Read: ReadIdentification,
		Size: ReflectSize(&Identification{}),
		Type: Both,
		Name: "Identification",
	})
}
