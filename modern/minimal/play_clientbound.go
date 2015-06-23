//go:generate protocol_generator -file=$GOFILE -direction=clientbound -state=play -package=minimal

package minimal

import "github.com/kurafuto/kyubu/packets"

// Packet ID: 0x01
type JoinGame struct {
	EntityID     int32
	Gamemode     uint8
	Dimension    int8
	Difficulty   uint8
	MaxPlayers   uint8
	LevelType    string
	ReducedDebug bool
}

// Packet ID: 0x02
type ServerMessage struct {
	Data     packets.Chat `as:"json"`
	Position int8
}

// Packet ID: 0x07
type Respawn struct {
	Dimension  int32
	Difficulty uint8
	Gamemode   uint8
	LevelType  string
}

// Packet ID: 0x38
type PlayerListItem struct {
	Action     packets.VarInt
	NumPlayers packets.VarInt

	Players []Player `length:"packets.VarInt"`
}

type Player struct {
	UUID           packets.UUID
	Name           string     `if:"t.Action == 0" noreplace:"true"`
	Properties     []Property `length:"packets.VarInt"`
	Gamemode       packets.VarInt
	Ping           packets.VarInt
	HasDisplayName bool
	DisplayName    packets.Chat `as:"json" if:".HasDisplayName"`

	// Action == 0 (add player)
	//  Name, Properties, Gamemode, Ping, HasDisplayName, DisplayName
	// Action == 1 (update gamemode)
	//  Gamemode
	// Action == 2 (update latency)
	//  Ping
	// Action == 3 (update display name)
	//  HasDisplayName, DisplayName
	// Action == 4 (remove player)
	//  [no fields]
}

type Property struct {
	Name      string
	Value     string
	Signed    bool
	Signature string `if:".Signed"`
}

// Packet ID: 0x3A
type ServerTabComplete struct {
	Matches []string `length:"packets.VarInt"`
}

// Packet ID: 0x3B
type ScoreboardObjective struct {
	Name  string
	Mode  uint8
	Value string `if:".Mode == 0 || .Mode == 2"`
	Type  string `if:".Mode == 0 || .Mode == 2"`
}

// Packet ID: 0x3C
type UpdateScore struct {
	Name          string
	Action        int8 // 0 = create/update, 1 = remove
	ObjectiveName string
	Value         packets.VarInt `if:".Action != 1"`
}

// Packet ID: 0x3D
type ShowScoreboard struct {
	Position int8
	Name     string
}

// Packet ID: 0x3E
type Teams struct {
	Name string
	// 0 = create, 1 = remove, 2 = update, 3 = add players, 4 = remove players
	Mode              int8
	Display           string   `if:".Mode == 0 || .Mode == 2"`
	Prefix            string   `if:".Mode == 0 || .Mode == 2"`
	Suffix            string   `if:".Mode == 0 || .Mode == 2"`
	FriendlyFire      int8     `if:".Mode == 0 || .Mode == 2"`
	NameTagVisibility string   `if:".Mode == 0 || .Mode == 2"`
	Color             int8     `if:".Mode == 0 || .Mode == 2"`
	Players           []string `length:"packets.VarInt" if:".Mode == 0 || .Mode == 3 || .Mode == 4"`
}

// Packet ID: 0x3F
type ServerPluginMessage struct {
	Channel string
	Data    []byte `length:"rest"`
}

// Packet ID: 0x40
type Disconnect struct {
	Reason packets.Chat `as:"json"`
}

// Packet ID: 0x46
type SetCompression struct {
	Threshold packets.VarInt
}
