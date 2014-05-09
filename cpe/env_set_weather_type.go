package cpe

import "github.com/sysr-q/kyubu/packets"

const EnvSetWeatherTypeSize = (packets.ByteSize + packets.ByteSize)

type EnvSetWeatherType struct {
	PacketId    byte
	WeatherType byte
}

func (p EnvSetWeatherType) Id() byte {
	return p.PacketId
}

func (p EnvSetWeatherType) Size() int {
	return EnvSetWeatherTypeSize
}

func (p EnvSetWeatherType) Bytes() []byte {
	return packets.ReflectBytes(p)
}

func (p EnvSetWeatherType) String() string {
	return "EnvWeatherType"
}

func ReadEnvSetWeatherType(b []byte) (packets.Packet, error) {
	var p EnvSetWeatherType
	err := packets.ReflectRead(b, &p)
	return &p, err
}

func NewEnvSetWeatherType(weatherType byte) (p *EnvSetWeatherType, err error) {
	p = &EnvSetWeatherType{
		PacketId:    0x1f,
		WeatherType: weatherType,
	}
	return
}

func init() {
	packets.MustRegister(&packets.PacketInfo{
		Id:   0x1f,
		Read: ReadEnvSetWeatherType,
		Size: EnvSetWeatherTypeSize,
		Type: packets.ServerOnly,
		Name: "Env Set Weather Type (CPE)",
	})
}
