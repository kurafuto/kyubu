package packets

type UpdateUserType struct {
	PacketId byte
	UserType byte
}

func (p UpdateUserType) Id() byte {
	return 0x0f
}

func (p UpdateUserType) Size() int {
	return ReflectSize(p)
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

func init() {
	MustRegister(&PacketInfo{
		Id:   0x0f,
		Read: ReadUpdateUserType,
		Size: ReflectSize(&UpdateUserType{}),
		Type: ServerOnly,
		Name: "Update User Type",
	})
}
