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

func (p *Identification) Id() byte {
	return p.PacketId
}

func (p *Identification) Size() int {
	return IdentificationSize
}

func (p *Identification) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteByte(p.ProtocolVersion)
	raw.WriteString(p.Name)
	raw.WriteString(p.KeyMotd)
	raw.WriteByte(p.UserType)
	return raw.Buffer.Bytes()
}

func ReadIdentification(b []byte) (Packet, error) {
	p := Identification{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if protocolVersion, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.ProtocolVersion = protocolVersion
	}
	if name, err := raw.ReadString(); err != nil {
		return nil, err
	} else {
		p.Name = name
	}
	if keyMotd, err := raw.ReadString(); err != nil {
		return nil, err
	} else {
		p.KeyMotd = keyMotd
	}
	if userType, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.UserType = userType
	}
	return &p, nil
}

func NewIdentification(name, key string) (p *Identification, err error) {
	if len(name) > StringSize || len(key) > StringSize {
		return nil, fmt.Errorf("kyubu/packets: cannot write over %d bytes in string", StringSize)
	}
	p = &Identification{
		PacketId:        0,
		ProtocolVersion: 7,
		Name:            name,
		KeyMotd:         key,
		UserType:        0,
	}
	return
}
