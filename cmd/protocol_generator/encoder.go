package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
)

type Encoder struct {
	p      packet
	prefix string
	buf    *bytes.Buffer

	tmpPrefix string
	tmpCount  int
}

func NewEncoder(p packet, b *bytes.Buffer) *Encoder {
	return &Encoder{
		p:      p,
		prefix: "t", // func (t *T) ...
		buf:    b,

		tmpPrefix: "tmp",
		tmpCount:  0,
	}
}

func (en *Encoder) T() string {
	en.tmpCount++
	return fmt.Sprintf("%s%d", en.tmpPrefix, en.tmpCount-1)
}

func (en *Encoder) Write() {
	en.writeStruct(en.p.spec, en.prefix)
}

func (en *Encoder) writeStruct(spec *ast.StructType, name string) {
	for _, field := range spec.Fields.List {
		// TODO: Parse field.Tag.Value -> reflect.StructTag
		tag := reflect.StructTag("")

		for _, n := range field.Names {
			en.writeType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
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

			fmt.Fprintf(en.buf, "// Encoding: %s (%s)\n", name, typeName)
		*/
	}
}

func (en *Encoder) writeType(e ast.Expr, name string, tag reflect.StructTag) {
	switch e := e.(type) {
	case *ast.StructType:
		en.writeStruct(e, name)
	case *ast.SelectorExpr:
		x := e.X.(*ast.Ident).Name
		s := e.Sel.Name
		en.writeField(x+"."+s, name, tag)
	case *ast.Ident:
		en.writeField(e.Name, name, tag)
	case *ast.ArrayType:
		// TODO
	default:
		fmt.Fprintf(en.buf, "// Unable to encode: %s (%T)\n", name, e)
	}
}

func (en *Encoder) writeField(t, name string, tag reflect.StructTag) {
	// TODO: For ints, unwrap binary.Write() trickery to reuse []byte tmp.
	switch t {
	case "bool":
		tmp := en.T() // byte value for bool
		fmt.Fprintf(en.buf, "%s := byte(0)\n", tmp)
		fmt.Fprintf(en.buf, "if %s {\n", name)
		fmt.Fprintf(en.buf, "\t%s = byte(1)\n", tmp)
		fmt.Fprintf(en.buf, "}\n")
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, tmp))
	case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, name))
	case "string":
		x := en.T() // []byte for varint
		b := en.T() // []byte(string)
		n := en.T() // num of varint bytes

		// TODO: Is making a new []byte every time efficient?
		fmt.Fprintf(en.buf, "%s := make([]byte, binary.MaxVarintLen32)\n", x)
		fmt.Fprintf(en.buf, "%s := []byte(%s)\n", b, name)
		fmt.Fprintf(en.buf, "%s := packets.PutVarint(%s, int64(len(%s)))\n", n, x, b)
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, b))
	case "packets.VarInt", "packets.VarLong":
		x := en.T() // []byte for varint
		n := en.T() // num of varint bytes

		// TODO: Is making a new []byte every time efficient?
		fmt.Fprintf(en.buf, "%s := make([]byte, binary.MaxVarintLen64)\n", x)
		fmt.Fprintf(en.buf, "%s := packets.PutVarint(%s, int64(%s))\n", n, x, name)
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
	case "packets.Position", "packets.Angle":
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, name))
	case "packets.UUID":
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s[:])", Endianness, name))

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
		fmt.Fprintf(en.buf, "// Unable to encode: %s (%s)\n", name, t)
	}

	fmt.Fprintf(en.buf, "\n")
}
