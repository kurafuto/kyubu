package cpe

import "github.com/sysr-q/kyubu/packets"

type HackControl struct {
	PacketId    byte
	Flying byte
	NoClip byte
	Speeding byte
	SpawnControl byte
	ThirdPersonView byte
	JumpHeight int16
}

func (p HackControl) Id() byte {
	return 0x20
}

func (p HackControl) Size() int {
	return packets.ReflectSize(p)
}

func (p HackControl) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p HackControl) String() string {
	return "HackControl"
}

func ReadHackControl(b []byte) (packets.Packet, error) {
	var p HackControl
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewHackControl(flying, noClip, speeding, spawnControl, thirdPerson byte, jumpHeight int16) (p *HackControl, err error) {
	p = &HackControl{
		PacketId:    0x20,
		Flying: flying,
		NoClip: noClip,
		Speeding: speeding,
		SpawnControl: spawnControl,
		ThirdPersonView: thirdPerson,
		JumpHeight: jumpHeight,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x20,
		Read: ReadHackControl,
		Size: packets.ReflectSize(&HackControl{}),
		Type: packets.ServerOnly,
		Name: "Hack Control (CPE)",
	})
}

