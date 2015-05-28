package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
)

type Decoder struct {
	p      packet
	prefix string
	buf    *bytes.Buffer

	tmpPrefix string
	tmpCount  int
}

func NewDecoder(p packet, b *bytes.Buffer) *Decoder {
	return &Decoder{
		p:      p,
		prefix: "t", // func (t *T) ...
		buf:    b,

		tmpPrefix: "tmp",
		tmpCount:  0,
	}
}

func (de *Decoder) T() string {
	de.tmpCount++
	return fmt.Sprintf("%s%d", de.tmpPrefix, de.tmpCount-1)
}

func (de *Decoder) Write() {
	de.writeStruct(de.p.spec, de.prefix)
}

func (de *Decoder) writeStruct(spec *ast.StructType, name string) {
	for _, field := range spec.Fields.List {
		// TODO: Parse field.Tag.Value -> reflect.StructTag
		tag := reflect.StructTag("")

		for _, n := range field.Names {
			de.writeType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
		}

		/*
			var typeName string

			if ide, ok := field.Type.(*ast.Ident); ok {
				typeName = ide.Name
			} else if selx, ok := field.Type.(*ast.SelectorExpr); ok {
				typeName = selx.X.(*ast.Ident).Name + "." + selx.Sel.Name
			} else if t, ok := field.Type.(*ast.ArrayType); ok {
				typeName = "[]" + t.Elt.(*ast.Ident).Name
			}

			fmt.Fprintf(de.buf, "// Decoding: %s (%s)\n", name, typeName)
		*/
	}
}

func (de *Decoder) writeType(e ast.Expr, name string, tag reflect.StructTag) {
	switch e := e.(type) {
	case *ast.StructType:
		de.writeStruct(e, name)
	case *ast.SelectorExpr:
		x := e.X.(*ast.Ident).Name
		s := e.Sel.Name
		de.writeField(x+"."+s, name, tag)
	case *ast.Ident:
		de.writeField(e.Name, name, tag)
	case *ast.ArrayType:
		// TODO
	default:
		fmt.Fprintf(de.buf, "// Unable to decode: %s (%T)\n", name, e)
	}
}

func (de *Decoder) writeField(t, name string, tag reflect.StructTag) {
	// TODO: For ints, unwrap binary.Read() trickery to reuse []byte tmp.
	switch t {
	case "bool":
		tmp := de.T() // byte value for bool
		fmt.Fprintf(de.buf, "var %s [1]byte\n", tmp)
		fmt.Fprintf(de.buf, "if _, err = rr.Read(%s[:1]); err != nil {\nreturn\n}\n", tmp)
		fmt.Fprintf(de.buf, "%s = %s[0] == 0x01\n", name, tmp)
	case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
		fmt.Fprintf(de.buf, errWrap("binary.Read(rr, %s, %s)", Endianness, name))
	case "string":
		imports["errors"] = struct{}{}

		n := de.T() // num of varint bytes
		b := de.T() // []byte to read into
		x := de.T() // num of bytes read

		fmt.Fprintf(de.buf, "%s, err := packets.ReadVarint(rr)\n", n)
		fmt.Fprintf(de.buf, "if err != nil {\nreturn\n}\n")
		fmt.Fprintf(de.buf, "%s := make([]byte, %s)\n", b, n)
		fmt.Fprintf(de.buf, "%s, err := rr.Read(%s)\n", x, b)
		fmt.Fprintf(de.buf, "if err != nil {\nreturn\n} ")
		fmt.Fprintf(de.buf, "else if int64(%s) != %s {\n", x, n)
		// TODO: delta := n - n1; Read(b, delta) ...
		fmt.Fprintf(de.buf, `return errors.New("didn't read enough bytes for string")`+"\n")
		fmt.Fprintf(de.buf, "}\n")
		fmt.Fprintf(de.buf, "%s = string(%s)\n", name, b)
	case "packets.VarInt":
		n := de.T()
		fmt.Fprintf(de.buf, "%s, err := packets.ReadVarint(rr)\n", n)
		fmt.Fprintf(de.buf, "if err != nil {\nreturn\n}\n")
		fmt.Fprintf(de.buf, "%s = packets.VarInt(%s)\n", name, n)
	case "packets.VarLong":
		n := de.T()
		fmt.Fprintf(de.buf, "%s, err := packets.ReadVarint(rr)\n", n)
		fmt.Fprintf(de.buf, "if err != nil {\nreturn\n}\n")
		fmt.Fprintf(de.buf, "%s = packets.VarLong(%s)\n", name, n)
	case "packets.Position", "packets.Angle":
		fmt.Fprintf(de.buf, errWrap("binary.Read(rr, %s, %s)", Endianness, name))
	case "packets.UUID":
		fmt.Fprintf(de.buf, errWrap("binary.Read(rr, %s, %s)", Endianness, name))

	case "packets.Chunk":
		fallthrough
	case "packets.Metadata":
		fallthrough
	case "packets.Slot":
		fallthrough
	case "packets.ObjectData":
		fallthrough
	case "packets.NBT":
		fallthrough
	default:
		fmt.Fprintf(de.buf, "// Unable to decode: %s (%s)\n", name, t)
	}

	fmt.Fprintf(de.buf, "\n")
}
