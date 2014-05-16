package packets

type OrientationUpdate struct {
	PacketId   byte
	PlayerId   int8
	Yaw, Pitch byte
}

func (p OrientationUpdate) Id() byte {
	return 0x0b
}

func (p OrientationUpdate) Size() int {
	return ReflectSize(p)
}

func (p OrientationUpdate) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadOrientationUpdate(b []byte) (Packet, error) {
	var p OrientationUpdate
	err := ReflectRead(b, &p)
	return &p, err
}

func NewOrientationUpdate(playerId int8, yaw, pitch byte) (p *OrientationUpdate, err error) {
	p = &OrientationUpdate{
		PacketId: 0x0b,
		PlayerId: playerId,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}

func init() {
	MustRegister(&PacketInfo{
		Id:   0x0b,
		Read: ReadOrientationUpdate,
		Size: ReflectSize(&OrientationUpdate{}),
		Type: ServerOnly,
		Name: "Orientation Update",
	})
}
