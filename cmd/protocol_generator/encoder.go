package main

import (
	"fmt"
	"go/ast"
)

type Encoder struct {
	p packet
	t string

	tmpCount int
}

func (en *Encoder) T() string {
	en.tmpCount++
	return fmt.Sprintf("tmp%d", en.tmpCount-1)
}

func (en *Encoder) Write() {
	for _, field := range en.p.spec.Fields.List {
		var typeName string

		if ide, ok := field.Type.(*ast.Ident); ok {
			typeName = ide.Name
		} else if selx, ok := field.Type.(*ast.SelectorExpr); ok {
			typeName = selx.X.(*ast.Ident).Name
			typeName = typeName + "." + selx.Sel.Name
		}

		// TODO: panic?
		name := field.Names[0].Name
		fmt.Fprintf(&buf, "// Encoding: %s (%s)\n", name, typeName)

		switch typeName {
		case "bool":
			tmp := en.T() // byte value for bool
			fmt.Fprintf(&buf, "%s := byte(0)\n", tmp)
			fmt.Fprintf(&buf, "if %s.%s {\n", en.t, name)
			fmt.Fprintf(&buf, "\t%s = byte(1)\n", tmp)
			fmt.Fprintf(&buf, "}\n")
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s)", Endianness, tmp))
		case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s.%s)", Endianness, en.t, name))
		case "string":
			x := en.T() // []byte for varint
			b := en.T() // []byte(string)
			n := en.T() // num of varint bytes

			// TODO: Is making a new []byte every time efficient?
			fmt.Fprintf(&buf, "%s := make([]byte, binary.MaxVarintLen32)\n", x)
			fmt.Fprintf(&buf, "%s := []byte(%s.%s)\n", b, en.t, name)
			fmt.Fprintf(&buf, "%s := packets.PutVarint(%s, int64(len(%s)))\n", n, x, b)
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s)", Endianness, b))
		case "packets.VarInt", "packets.VarLong":
			x := en.T() // []byte for varint
			n := en.T() // num of varint bytes

			// TODO: Is making a new []byte every time efficient?
			fmt.Fprintf(&buf, "%s := make([]byte, binary.MaxVarintLen64)\n", x)
			fmt.Fprintf(&buf, "%s := packets.PutVarint(%s, int64(%s.%s))\n", n, x, en.t, name)
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
		case "packets.Chunk":
		case "packets.Metadata":
		case "packets.Slot":
		case "packets.ObjectData":
		case "packets.NBT":
		case "packets.Position", "packets.Angle":
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s.%s)", Endianness, en.t, name))
		case "packets.UUID":
			fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s.%s[:])", Endianness, en.t, name))
		default:
			fmt.Fprintf(&buf, "// Unable to encode: %s.%s (%s)\n", en.t, name, typeName)
		}

		fmt.Fprintf(&buf, "\n")
	}
}
