package kyubu

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
