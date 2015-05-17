package classic

import "github.com/sysr-q/kyubu/packets"

type UpdateUserType struct {
	PacketId byte
	UserType byte
}

func (p UpdateUserType) Id() byte {
	return 0x0f
}

func (p UpdateUserType) Size() int {
	return packets.ReflectSize(p)
}

func (p UpdateUserType) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadUpdateUserType(b []byte) (packets.Packet, error) {
	var p UpdateUserType
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewUpdateUserType(userType byte) (p *UpdateUserType, err error) {
	p = &UpdateUserType{
		PacketId:      0x0f,
		UserType: userType,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x0f,
		Read:      ReadUpdateUserType,
		Size:      packets.ReflectSize(&UpdateUserType{}),
		Direction: packets.Anomalous,
		Name:      "Update User Type",
	})
}
