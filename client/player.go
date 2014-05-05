package client

import "github.com/sysr-q/kyubu/packets"

type Player struct {
	Id              int8
	Name            string
	Location, Spawn *Loc
	Op              bool
}

func (p *Player) Teleport(loc *Loc) {
	p.Location.X = loc.X
	p.Location.Y = loc.Y
	p.Location.Z = loc.Z
	p.Location.Yaw = loc.Yaw
	p.Location.Pitch = loc.Pitch
}

type Loc struct {
	X, Y, Z    int16
	Yaw, Pitch byte
}

func POLoc(p *packets.PositionOrientation) *Loc {
	return &Loc{
		X:     p.X,
		Y:     p.Y,
		Z:     p.Z,
		Yaw:   p.Yaw,
		Pitch: p.Pitch,
	}
}
