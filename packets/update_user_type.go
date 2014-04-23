package packets

const UpdateUserTypeSize = (ByteSize + // Packet ID
	ByteSize) // User type

type UpdateUserType struct {
	PacketId byte
	UserType byte
}

func (p *UpdateUserType) Id() byte {
	return p.PacketId
}

func (p *UpdateUserType) Size() int {
	return UpdateUserTypeSize
}

func (p *UpdateUserType) Bytes() []byte {
	raw := NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteByte(p.UserType)
	return raw.Buffer.Bytes()
}

func ReadUpdateUserType(b []byte) (Packet, error) {
	p := UpdateUserType{}
	raw := NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if userType, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.UserType = userType
	}
	return &p, nil
}

func NewUpdateUserType(userType byte) (p *UpdateUserType, err error) {
	p = &UpdateUserType{
		PacketId: 12,
		UserType: userType,
	}
	return
}
