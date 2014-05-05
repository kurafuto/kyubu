package packets

const UpdateUserTypeSize = (ByteSize + // Packet ID
	ByteSize) // User type

type UpdateUserType struct {
	PacketId byte
	UserType byte
}

func (p UpdateUserType) Id() byte {
	return p.PacketId
}

func (p UpdateUserType) Size() int {
	return UpdateUserTypeSize
}

func (p UpdateUserType) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadUpdateUserType(b []byte) (Packet, error) {
	var p UpdateUserType
	err := ReflectRead(b, &p)
	return &p, err
}

func NewUpdateUserType(userType byte) (p *UpdateUserType, err error) {
	p = &UpdateUserType{
		PacketId: 0x0f,
		UserType: userType,
	}
	return
}
