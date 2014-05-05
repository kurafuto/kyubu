package packets

import "fmt"

const IdentificationSize = (ByteSize + // Packet ID
	ByteSize + // Protocol Version
	StringSize + // Client: Username, Server: Server name
	StringSize + // Client: Verification key, Server: Server MOTD
	ByteSize) // Client: unused, Server: User type

type Identification struct {
	PacketId        byte
	ProtocolVersion byte
	Name            string
	KeyMotd         string
	UserType        byte
}

func (p Identification) Id() byte {
	return p.PacketId
}

func (p Identification) Size() int {
	return IdentificationSize
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
