package classic

import "github.com/sysr-q/kyubu/packets"

type LevelFinalize struct {
	PacketId byte
	X, Y, Z  int16
}

func (p LevelFinalize) Id() byte {
	return 0x04
}

func (p LevelFinalize) Size() int {
	return packets.ReflectSize(p)
}

func (p LevelFinalize) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func ReadLevelFinalize(b []byte) (packets.Packet, error) {
	var p LevelFinalize
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewLevelFinalize(x, y, z int16) (p *LevelFinalize, err error) {
	p = &LevelFinalize{
		PacketId: 0x04,
		X:        x,
		Y:        y,
		Z:        z,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:        0x04,
		Read:      ReadLevelFinalize,
		Size:      packets.ReflectSize(&LevelFinalize{}),
		Direction: packets.Anomalous,
		Name:      "Level Finalize",
	})
}
