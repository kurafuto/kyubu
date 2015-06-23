//go:generate protocol_generator -file=$GOFILE -direction=serverbound -state=status -package=modern

package modern

import "github.com/kurafuto/kyubu/packets"

// Packet ID: 0x00
type StatusResponse struct {
	Status StatusReply `as:"json"`
}

// Packet ID: 0x01
type StatusPong struct {
	Time int64
}

// Lifted from thinkofdeath/steven.

// StatusReply is the reply retrieved from a server when pinging it.
type StatusReply struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int            `json:"max"`
		Online int            `json:"online"`
		Sample []StatusPlayer `json:"sample,omitempty"`
	} `json:"players"`
	Description packets.Chat `json:"description"`
	Favicon     string       `json:"favicon"`
}

// StatusPlayer is one of the sample players in a StatusReply
type StatusPlayer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
