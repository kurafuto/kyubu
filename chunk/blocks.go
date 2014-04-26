package chunk

type BlockId byte

const (
	Air = BlockId(iota)
	Stone
	Grass
	Dirt
	Cobblestone
	WoodenPlank
	Sapling
	Bedrock
	Water
	StationaryWater
	Lava
	StationaryLava
	Sand
	Gravel
	GoldOre
	IronOre
	CoalOre
	Wood
	Leaves
	Sponge
	Glass
	RedCloth
	OrangeCloth
	YellowCloth
	LimeCloth
	GreenCloth
	AquaGreenCloth
	CyanCloth
	BlueCloth
	PurpleCloth
	IndigoCloth
	VioletCloth
	MagentaCloth
	PinkCloth
	BlackCloth
	GrayCloth
	WhiteCloth
	Dandelion
	Rose
	BrownMushroom
	RedMushroom
	GoldBlock
	IronBlock
	DoubleSlab
	Slab
	BrickBlock
	TNT
	Bookshelf
	MossStone
	Obsidian
)

var Blocks = map[BlockId]string{
	0x00: "Air",
	0x01: "Stone",
	0x02: "Grass",
	0x03: "Dirt",
	0x04: "Cobblestone",
	0x05: "Wooden Plank",
	0x06: "Sapling",
	0x07: "Bedrock",
	0x08: "Water",
	0x09: "Stationary Water",
	0x0a: "Lava",
	0x0b: "Stationary Lava",
	0x0c: "Sand",
	0x0d: "Gravel",
	0x0e: "Gold Ore",
	0x0f: "Iron Ore",
	0x10: "Coal Ore",
	0x11: "Wood",
	0x12: "Leaves",
	0x13: "Sponge",
	0x14: "Glass",
	0x15: "Red Cloth",
	0x16: "Orange Cloth",
	0x17: "Yellow Cloth",
	0x18: "Lime Cloth",
	0x19: "Green Cloth",
	0x1a: "Aqua Green Cloth",
	0x1b: "Cyan Cloth",
	0x1c: "Blue Cloth",
	0x1d: "Purple Cloth",
	0x1e: "Indigo Cloth",
	0x1f: "Violet Cloth",
	0x20: "Magenta Cloth",
	0x21: "Pink Cloth",
	0x22: "Black Cloth",
	0x23: "Gray Cloth",
	0x24: "White Cloth",
	0x25: "Dandelion",
	0x26: "Rose",
	0x27: "Brown Mushroom",
	0x28: "Red Mushroom",
	0x29: "Gold Block",
	0x2a: "Iron Block",
	0x2b: "Double Slab",
	0x2c: "Slab",
	0x2d: "Brick Block",
	0x2e: "TNT",
	0x2f: "Bookshelf",
	0x30: "Moss Stone",
	0x31: "Obsidian",
}
