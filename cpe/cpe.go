// Package cpe implements all of the Classic Protocol Extension packets in the
// official CPE spec.
package cpe

import "github.com/sysr-q/kyubu/packets"

// ExtPacket is, in all aspects, the same as kyubu/packets.Packet, but has an
// extra String() method, which will tell you which extension set the packet is from
//
// All packets in kyubu/cpe implement this interface.
type ExtPacket interface {
	packets.Packet
	String() string // The CPE extension name this packet is related to.
}
