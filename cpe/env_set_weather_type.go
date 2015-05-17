package cpe

import "github.com/sysr-q/kyubu/packets"

type EnvSetWeatherType struct {
	PacketId    byte
	WeatherType byte
}

func (p EnvSetWeatherType) Id() byte {
	return 0x1f
}

func (p EnvSetWeatherType) Size() int {
	return packets.ReflectSize(p)
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
	packets.Register(&packets.PacketInfo{
		Id:        0x1f,
		Read:      ReadEnvSetWeatherType,
		Size:      packets.ReflectSize(&EnvSetWeatherType{}),
		Direction: packets.Anomalous,
		Name:      "Env Set Weather Type (CPE)",
	})
}
