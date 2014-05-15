package cpe

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

const ChangeModelSize = (packets.ByteSize + packets.ByteSize + packets.StringSize)

type ChangeModel struct {
	PacketId  byte
	EntityId  byte
	ModelName string
}

func (p ChangeModel) Id() byte {
	return 0x1d
}

func (p ChangeModel) Size() int {
	return ChangeModelSize
}

func (p ChangeModel) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p ChangeModel) String() string {
	return "ChangeModel"
}

func ReadChangeModel(b []byte) (packets.Packet, error) {
	var p ChangeModel
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewChangeModel(entityId byte, modelName string) (p *ChangeModel, err error) {
	if len(modelName) > packets.StringSize {
		return nil, fmt.Errorf("kyubu/cpe: cannot write over %d bytes in string", packets.StringSize)
	}

	p = &ChangeModel{
		PacketId:  0x1d,
		EntityId:  entityId,
		ModelName: modelName,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x1d,
		Read: ReadChangeModel,
		Size: ChangeModelSize,
		Type: packets.ServerOnly,
		Name: "Change Model (CPE)",
	})
}
