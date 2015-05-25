package main

import (
	"fmt"
	"go/ast"
)

type Decoder struct {
	p packet
	t string

	tmpPrefix string
	tmpCount  int
}

func (de *Decoder) T() string {
	de.tmpCount++
	return fmt.Sprintf("%s%d", de.tmpPrefix, de.tmpCount-1)
}

func (de *Decoder) Write() {
	for _, field := range de.p.spec.Fields.List {
		var typeName string

		if ide, ok := field.Type.(*ast.Ident); ok {
			typeName = ide.Name
		} else if selx, ok := field.Type.(*ast.SelectorExpr); ok {
			typeName = selx.X.(*ast.Ident).Name
			typeName = typeName + "." + selx.Sel.Name
		}

		// TODO: panic?
		name := field.Names[0].Name
		fmt.Fprintf(&buf, "// Decoding: %s (%s)\n", name, typeName)

		switch typeName {
		case "bool":
			tmp := de.T() // byte value for bool
			fmt.Fprintf(&buf, "var %s [1]byte\n", tmp)
			fmt.Fprintf(&buf, "if _, err = rr.Read(%s[:1]); err != nil {\nreturn\n}\n", tmp)
			fmt.Fprintf(&buf, "%s.%s = %s[0] == 0x01\n", de.t, name, tmp)
		case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
			fmt.Fprintf(&buf, errWrap("binary.Read(rr, %s, %s.%s)", Endianness, de.t, name))
		case "string":
			imports["errors"] = struct{}{}

			n := de.T() // num of varint bytes
			b := de.T() // []byte to read into
			x := de.T() // num of bytes read

			fmt.Fprintf(&buf, "%s, err := packets.ReadVarint(rr)\n", n)
			fmt.Fprintf(&buf, "if err != nil {\nreturn\n}\n")
			fmt.Fprintf(&buf, "%s := make([]byte, %s)\n", b, n)
			fmt.Fprintf(&buf, "%s, err := rr.Read(%s)\n", x, b)
			fmt.Fprintf(&buf, "if err != nil {\nreturn\n} ")
			fmt.Fprintf(&buf, "else if int64(%s) != %s {\n", x, n)
			// TODO: delta := n - n1; Read(b, delta) ...
			fmt.Fprintf(&buf, `return errors.New("didn't read enough bytes for string")`+"\n")
			fmt.Fprintf(&buf, "}\n")
			fmt.Fprintf(&buf, "%s.%s = string(%s)\n", de.t, name, b)
		case "packets.VarInt":
			n := de.T()
			fmt.Fprintf(&buf, "%s, err := packets.ReadVarint(rr)\n", n)
			fmt.Fprintf(&buf, "if err != nil {\nreturn\n}\n")
			fmt.Fprintf(&buf, "%s.%s = packets.VarInt(%s)\n", de.t, name, n)
		case "packets.VarLong":
			n := de.T()
			fmt.Fprintf(&buf, "%s, err := packets.ReadVarint(rr)\n", n)
			fmt.Fprintf(&buf, "if err != nil {\nreturn\n}\n")
			fmt.Fprintf(&buf, "%s.%s = packets.VarLong(%s)\n", de.t, name, n)
		case "packets.Chunk":
		case "packets.Metadata":
		case "packets.Slot":
		case "packets.ObjectData":
		case "packets.NBT":
		case "packets.Position", "packets.Angle":
			fmt.Fprintf(&buf, errWrap("binary.Read(rr, %s, %s.%s)", Endianness, de.t, name))
		case "packets.UUID":
			fmt.Fprintf(&buf, errWrap("binary.Read(rr, %s, %s.%s)", Endianness, de.t, name))
		default:
			fmt.Fprintf(&buf, "// Unable to decode: %s.%s (%s)\n", de.t, name, typeName)
		}

		fmt.Fprintf(&buf, "\n")
	}
}
