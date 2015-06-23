package minimal

import (
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/kurafuto/kyubu/packets"
)

func init() {
	//// Play: ClientBound
	// Packet ID: 0x04: Entity Equipment
	packets.Register(packets.Play, packets.ClientBound, 0x04, func() packets.Packet { return &GenericPacket{id: 0x04} })
	// Packet ID: 0x05: Spawn Position
	packets.Register(packets.Play, packets.ClientBound, 0x05, func() packets.Packet { return &GenericPacket{id: 0x05} })
	// Packet ID: 0x06: Update Health
	packets.Register(packets.Play, packets.ClientBound, 0x06, func() packets.Packet { return &GenericPacket{id: 0x06} })
	// Packet ID: 0x07: Respawn
	packets.Register(packets.Play, packets.ClientBound, 0x07, func() packets.Packet { return &GenericPacket{id: 0x07} })
	// Packet ID: 0x08: Player Position And Look
	packets.Register(packets.Play, packets.ClientBound, 0x08, func() packets.Packet { return &GenericPacket{id: 0x08} })
	// Packet ID: 0x09: Held Item Change
	packets.Register(packets.Play, packets.ClientBound, 0x09, func() packets.Packet { return &GenericPacket{id: 0x09} })
	// Packet ID: 0x0A: Use Bed
	packets.Register(packets.Play, packets.ClientBound, 0x0A, func() packets.Packet { return &GenericPacket{id: 0x0A} })
	// Packet ID: 0x0B: Animation
	packets.Register(packets.Play, packets.ClientBound, 0x0B, func() packets.Packet { return &GenericPacket{id: 0x0B} })
	// Packet ID: 0x0C: Spawn Player
	packets.Register(packets.Play, packets.ClientBound, 0x0C, func() packets.Packet { return &GenericPacket{id: 0x0C} })
	// Packet ID: 0x0D: Collect Item
	packets.Register(packets.Play, packets.ClientBound, 0x0D, func() packets.Packet { return &GenericPacket{id: 0x0D} })
	// Packet ID: 0x0E: Spawn Object
	packets.Register(packets.Play, packets.ClientBound, 0x0E, func() packets.Packet { return &GenericPacket{id: 0x0E} })
	// Packet ID: 0x0F: Spawn Mob
	packets.Register(packets.Play, packets.ClientBound, 0x0F, func() packets.Packet { return &GenericPacket{id: 0x0F} })
	// Packet ID: 0x10: Spawn Painting
	packets.Register(packets.Play, packets.ClientBound, 0x10, func() packets.Packet { return &GenericPacket{id: 0x10} })
	// Packet ID: 0x11: Spawn Experience Orb
	packets.Register(packets.Play, packets.ClientBound, 0x11, func() packets.Packet { return &GenericPacket{id: 0x11} })
	// Packet ID: 0x12: Entity Velocity
	packets.Register(packets.Play, packets.ClientBound, 0x12, func() packets.Packet { return &GenericPacket{id: 0x12} })
	// Packet ID: 0x13: Destroy Entities
	packets.Register(packets.Play, packets.ClientBound, 0x13, func() packets.Packet { return &GenericPacket{id: 0x13} })
	// Packet ID: 0x14: Entity
	packets.Register(packets.Play, packets.ClientBound, 0x14, func() packets.Packet { return &GenericPacket{id: 0x14} })
	// Packet ID: 0x15: Entity Relative Move
	packets.Register(packets.Play, packets.ClientBound, 0x15, func() packets.Packet { return &GenericPacket{id: 0x15} })
	// Packet ID: 0x16: Entity Look
	packets.Register(packets.Play, packets.ClientBound, 0x16, func() packets.Packet { return &GenericPacket{id: 0x16} })
	// Packet ID: 0x17: Entity Look And Relative Move
	packets.Register(packets.Play, packets.ClientBound, 0x17, func() packets.Packet { return &GenericPacket{id: 0x17} })
	// Packet ID: 0x18: Entity Teleport
	packets.Register(packets.Play, packets.ClientBound, 0x18, func() packets.Packet { return &GenericPacket{id: 0x18} })
	// Packet ID: 0x19: Entity Head Look
	packets.Register(packets.Play, packets.ClientBound, 0x19, func() packets.Packet { return &GenericPacket{id: 0x19} })
	// Packet ID: 0x1A: Entity Status
	packets.Register(packets.Play, packets.ClientBound, 0x1A, func() packets.Packet { return &GenericPacket{id: 0x1A} })
	// Packet ID: 0x1B: Attach Entity
	packets.Register(packets.Play, packets.ClientBound, 0x1B, func() packets.Packet { return &GenericPacket{id: 0x1B} })
	// Packet ID: 0x1C: Entity Metadata
	packets.Register(packets.Play, packets.ClientBound, 0x1C, func() packets.Packet { return &GenericPacket{id: 0x1C} })
	// Packet ID: 0x1D: Entity Effect
	packets.Register(packets.Play, packets.ClientBound, 0x1D, func() packets.Packet { return &GenericPacket{id: 0x1D} })
	// Packet ID: 0x1E: Remove Entity Effect
	packets.Register(packets.Play, packets.ClientBound, 0x1E, func() packets.Packet { return &GenericPacket{id: 0x1E} })
	// Packet ID: 0x1F: Set Experience
	packets.Register(packets.Play, packets.ClientBound, 0x1F, func() packets.Packet { return &GenericPacket{id: 0x1F} })
	// Packet ID: 0x20: Entity Properties
	packets.Register(packets.Play, packets.ClientBound, 0x20, func() packets.Packet { return &GenericPacket{id: 0x20} })
	// Packet ID: 0x21: Chunk Data
	packets.Register(packets.Play, packets.ClientBound, 0x21, func() packets.Packet { return &GenericPacket{id: 0x21} })
	// Packet ID: 0x22: Multi Block Change
	packets.Register(packets.Play, packets.ClientBound, 0x22, func() packets.Packet { return &GenericPacket{id: 0x22} })
	// Packet ID: 0x23: Block Change
	packets.Register(packets.Play, packets.ClientBound, 0x23, func() packets.Packet { return &GenericPacket{id: 0x23} })
	// Packet ID: 0x24: Block Action
	packets.Register(packets.Play, packets.ClientBound, 0x24, func() packets.Packet { return &GenericPacket{id: 0x24} })
	// Packet ID: 0x25: Block Break Animation
	packets.Register(packets.Play, packets.ClientBound, 0x25, func() packets.Packet { return &GenericPacket{id: 0x25} })
	// Packet ID: 0x26: Map Change Bulk
	packets.Register(packets.Play, packets.ClientBound, 0x26, func() packets.Packet { return &GenericPacket{id: 0x26} })
	// Packet ID: 0x27: Explosion
	packets.Register(packets.Play, packets.ClientBound, 0x27, func() packets.Packet { return &GenericPacket{id: 0x27} })
	// Packet ID: 0x28: Effect
	packets.Register(packets.Play, packets.ClientBound, 0x28, func() packets.Packet { return &GenericPacket{id: 0x28} })
	// Packet ID: 0x29: Sound Effect
	packets.Register(packets.Play, packets.ClientBound, 0x29, func() packets.Packet { return &GenericPacket{id: 0x29} })
	// Packet ID: 0x2A: Particle
	packets.Register(packets.Play, packets.ClientBound, 0x2A, func() packets.Packet { return &GenericPacket{id: 0x2A} })
	// Packet ID: 0x2B: Change Game State
	packets.Register(packets.Play, packets.ClientBound, 0x2B, func() packets.Packet { return &GenericPacket{id: 0x2B} })
	// Packet ID: 0x2C: Spawn Global Entity
	packets.Register(packets.Play, packets.ClientBound, 0x2C, func() packets.Packet { return &GenericPacket{id: 0x2C} })
	// Packet ID: 0x2D: Open Window
	packets.Register(packets.Play, packets.ClientBound, 0x2D, func() packets.Packet { return &GenericPacket{id: 0x2D} })
	// Packet ID: 0x2E: Close Window
	packets.Register(packets.Play, packets.ClientBound, 0x2E, func() packets.Packet { return &GenericPacket{id: 0x2E} })
	// Packet ID: 0x2F: Set Slot
	packets.Register(packets.Play, packets.ClientBound, 0x2F, func() packets.Packet { return &GenericPacket{id: 0x2F} })
	// Packet ID: 0x30: Window Items
	packets.Register(packets.Play, packets.ClientBound, 0x30, func() packets.Packet { return &GenericPacket{id: 0x30} })
	// Packet ID: 0x31: Window Property
	packets.Register(packets.Play, packets.ClientBound, 0x31, func() packets.Packet { return &GenericPacket{id: 0x31} })
	// Packet ID: 0x32: Confirm Transaction
	packets.Register(packets.Play, packets.ClientBound, 0x32, func() packets.Packet { return &GenericPacket{id: 0x32} })
	// Packet ID: 0x33: Update Sign
	packets.Register(packets.Play, packets.ClientBound, 0x33, func() packets.Packet { return &GenericPacket{id: 0x33} })
	// Packet ID: 0x34: Maps
	packets.Register(packets.Play, packets.ClientBound, 0x34, func() packets.Packet { return &GenericPacket{id: 0x34} })
	// Packet ID: 0x35: Update Block Entity
	packets.Register(packets.Play, packets.ClientBound, 0x35, func() packets.Packet { return &GenericPacket{id: 0x35} })
	// Packet ID: 0x36: Sign Editor Open
	packets.Register(packets.Play, packets.ClientBound, 0x36, func() packets.Packet { return &GenericPacket{id: 0x36} })
	// Packet ID: 0x37: Statistics
	packets.Register(packets.Play, packets.ClientBound, 0x37, func() packets.Packet { return &GenericPacket{id: 0x37} })
	// Packet ID: 0x39: Player Abilities
	packets.Register(packets.Play, packets.ClientBound, 0x39, func() packets.Packet { return &GenericPacket{id: 0x39} })
	// Packet ID: 0x41: Server Difficulty
	packets.Register(packets.Play, packets.ClientBound, 0x41, func() packets.Packet { return &GenericPacket{id: 0x41} })
	// Packet ID: 0x42: Combat Event
	packets.Register(packets.Play, packets.ClientBound, 0x42, func() packets.Packet { return &GenericPacket{id: 0x42} })
	// Packet ID: 0x43: Camera
	packets.Register(packets.Play, packets.ClientBound, 0x43, func() packets.Packet { return &GenericPacket{id: 0x43} })
	// Packet ID: 0x44: World Border
	packets.Register(packets.Play, packets.ClientBound, 0x44, func() packets.Packet { return &GenericPacket{id: 0x44} })
	// Packet ID: 0x45: Title
	packets.Register(packets.Play, packets.ClientBound, 0x45, func() packets.Packet { return &GenericPacket{id: 0x45} })
	// Packet ID: 0x47: Player List Header/Footer
	packets.Register(packets.Play, packets.ClientBound, 0x47, func() packets.Packet { return &GenericPacket{id: 0x47} })
	// Packet ID: 0x48: Resource Pack Send
	packets.Register(packets.Play, packets.ClientBound, 0x48, func() packets.Packet { return &GenericPacket{id: 0x48} })
	// Packet ID: 0x49: Update Entity NBT
	packets.Register(packets.Play, packets.ClientBound, 0x49, func() packets.Packet { return &GenericPacket{id: 0x49} })

	//// Play: ServerBound
	// Packet ID: 0x00: Keep Alive
	packets.Register(packets.Play, packets.ServerBound, 0x00, func() packets.Packet { return &GenericPacket{id: 0x00} })
	// Packet ID: 0x02: Use Entity
	packets.Register(packets.Play, packets.ServerBound, 0x02, func() packets.Packet { return &GenericPacket{id: 0x02} })
	// Packet ID: 0x03: Player
	packets.Register(packets.Play, packets.ServerBound, 0x03, func() packets.Packet { return &GenericPacket{id: 0x03} })
	// Packet ID: 0x04: Player Position
	packets.Register(packets.Play, packets.ServerBound, 0x04, func() packets.Packet { return &GenericPacket{id: 0x04} })
	// Packet ID: 0x05: Player Look
	packets.Register(packets.Play, packets.ServerBound, 0x05, func() packets.Packet { return &GenericPacket{id: 0x05} })
	// Packet ID: 0x06: Player Position And Look
	packets.Register(packets.Play, packets.ServerBound, 0x06, func() packets.Packet { return &GenericPacket{id: 0x06} })
	// Packet ID: 0x07: Player Digging
	packets.Register(packets.Play, packets.ServerBound, 0x07, func() packets.Packet { return &GenericPacket{id: 0x07} })
	// Packet ID: 0x08: Player Block Placement
	packets.Register(packets.Play, packets.ServerBound, 0x08, func() packets.Packet { return &GenericPacket{id: 0x08} })
	// Packet ID: 0x09: Held Item Change
	packets.Register(packets.Play, packets.ServerBound, 0x09, func() packets.Packet { return &GenericPacket{id: 0x09} })
	// Packet ID: 0x0A: Animation
	packets.Register(packets.Play, packets.ServerBound, 0x0A, func() packets.Packet { return &GenericPacket{id: 0x0A} })
	// Packet ID: 0x0B: Entity Action
	packets.Register(packets.Play, packets.ServerBound, 0x0B, func() packets.Packet { return &GenericPacket{id: 0x0B} })
	// Packet ID: 0x0C: Steer Vehicle
	packets.Register(packets.Play, packets.ServerBound, 0x0C, func() packets.Packet { return &GenericPacket{id: 0x0C} })
	// Packet ID: 0x0D: Close Window
	packets.Register(packets.Play, packets.ServerBound, 0x0D, func() packets.Packet { return &GenericPacket{id: 0x0D} })
	// Packet ID: 0x0E: Click Window
	packets.Register(packets.Play, packets.ServerBound, 0x0E, func() packets.Packet { return &GenericPacket{id: 0x0E} })
	// Packet ID: 0x0F: Confirm Transaction
	packets.Register(packets.Play, packets.ServerBound, 0x0F, func() packets.Packet { return &GenericPacket{id: 0x0F} })
	// Packet ID: 0x10: Creative Inventory Action
	packets.Register(packets.Play, packets.ServerBound, 0x10, func() packets.Packet { return &GenericPacket{id: 0x10} })
	// Packet ID: 0x11: Enchant Item
	packets.Register(packets.Play, packets.ServerBound, 0x11, func() packets.Packet { return &GenericPacket{id: 0x11} })
	// Packet ID: 0x12: Update Sign
	packets.Register(packets.Play, packets.ServerBound, 0x12, func() packets.Packet { return &GenericPacket{id: 0x12} })
	// Packet ID: 0x13: Player Abilities
	packets.Register(packets.Play, packets.ServerBound, 0x13, func() packets.Packet { return &GenericPacket{id: 0x13} })
	// Packet ID: 0x15: Client Settings
	packets.Register(packets.Play, packets.ServerBound, 0x15, func() packets.Packet { return &GenericPacket{id: 0x15} })
	// Packet ID: 0x18: Spectate
	packets.Register(packets.Play, packets.ServerBound, 0x18, func() packets.Packet { return &GenericPacket{id: 0x18} })
	// Packet ID: 0x19: Resource Pack Status
	packets.Register(packets.Play, packets.ServerBound, 0x19, func() packets.Packet { return &GenericPacket{id: 0x19} })
}

type GenericPacket struct {
	id   byte
	Data []byte
}

func (t *GenericPacket) Id() byte {
	return t.id
}

func (t *GenericPacket) Encode(ww io.Writer) error {
	return binary.Write(ww, binary.BigEndian, t.Data)
}

func (t *GenericPacket) Decode(rr io.Reader) error {
	b, err := ioutil.ReadAll(rr)
	if err != nil {
		return err
	}

	t.Data = b
	return nil
}
