package packets

type PositionOrientation struct {
	PacketId   byte
	PlayerId   int8
	X, Y, Z    int16
	Yaw, Pitch byte
}

func (p PositionOrientation) Id() byte {
	return 0x08
}

func (p PositionOrientation) Size() int {
	return ReflectSize(p)
}

func (p PositionOrientation) Bytes() []byte {
	return ReflectBytes(p)
}

func ReadPositionOrientation(b []byte) (Packet, error) {
	var p PositionOrientation
	err := ReflectRead(b, &p)
	return &p, err
}

func NewPositionOrientation(playerId int8, x, y, z int16, yaw, pitch byte) (p *PositionOrientation, err error) {
	p = &PositionOrientation{
		PacketId: 0x08,
		PlayerId: playerId,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      yaw,
		Pitch:    pitch,
	}
	return
}

func init() {
	MustRegister(&PacketInfo{
		Id:   0x08,
		Read: ReadPositionOrientation,
		Size: ReflectSize(&PositionOrientation{}),
		Type: Both,
		Name: "Position/Orientation",
	})
}
