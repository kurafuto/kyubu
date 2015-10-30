package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
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
		tag := reflect.StructTag("")
		if field.Tag != nil {
			tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		}

		var c string
		if tag.Get("if") != "" && tag.Get("noreplace") != "true" {
			c = strings.Replace(tag.Get("if"), ".", name+".", -1)
		}

		if c != "" {
			fmt.Fprintf(en.buf, "if %s {\n", c)
		}

		for _, n := range field.Names {
			en.writeType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
		}

		if c != "" {
			fmt.Fprint(en.buf, "}\n")
		}
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
		lT := tag.Get("length") // length type
		if lT != "rest" && lT != "" {
			en.writeField(lT, fmt.Sprintf("len(%s)", name), "")
		}

		if i, ok := e.Elt.(*ast.Ident); ok && (i.Name == "byte" || i.Name == "uint8") {
			fmt.Fprintf(en.buf, "if _, err = ww.Write(%s); err != nil { return err }\n", name)
			return
		}

		iv := en.T()
		fmt.Fprintf(en.buf, "for %s := range %s {\n", iv, name)

		// Sneaky trick so we can use "unknown-ish" structs ([]T) and still get at their fields.
		sName := e.Elt.(*ast.Ident).Name
		if xx, ok := nonPackets[sName]; ok {
			en.writeType(xx, fmt.Sprintf("%s[%s]", name, iv), tag)
		} else if sName == "string" {
			en.writeField(sName, fmt.Sprintf("%s[%s]", name, iv), tag)
		} else {
			fmt.Fprintf(en.buf, "// Can't find supporting struct %s for %s.%s\n", sName, en.p.name, name)
		}

		fmt.Fprint(en.buf, "}\n")
	default:
		fmt.Fprintf(en.buf, "// Unable to encode: %s (%T)\n", name, e)
	}
}

func (en *Encoder) writeField(t, name string, tag reflect.StructTag) {
	as := tag.Get("as")
	if as != "" {
		switch as {
		case "json":
			imports["encoding/json"] = struct{}{}
			t := en.T()
			fmt.Fprintf(en.buf, `var %[1]s []byte
			if %[1]s, err = json.Marshal(&%[2]s); err != nil { return err }
			if err = packets.WriteString(ww, string(%[1]s)); err != nil { return err }
			`, t, name)
		default:
			fmt.Fprintf(en.buf, "// Can't 'as' %s\n", as)
		}
		return
	}

	// TODO: For ints, unwrap binary.Write() trickery to reuse []byte tmp.
	switch t {
	case "bool":
		/*
			tmp := en.T() // byte value for bool
			fmt.Fprintf(en.buf, "%s := byte(0)\n", tmp)
			fmt.Fprintf(en.buf, "if %s {\n", name)
			fmt.Fprintf(en.buf, "\t%s = byte(1)\n", tmp)
			fmt.Fprintf(en.buf, "}\n")
			fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, tmp))
		*/
		fmt.Fprintf(en.buf, "if err = packets.WriteBool(ww, %s); err != nil { return err }", name)
	case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, name))
	case "string":
		/*
			x := en.T() // []byte for varint
			b := en.T() // []byte(string)
			n := en.T() // num of varint bytes

			// TODO: Is making a new []byte every time efficient?
			fmt.Fprintf(en.buf, "%s := make([]byte, binary.MaxVarintLen32)\n", x)
			fmt.Fprintf(en.buf, "%s := []byte(%s)\n", b, name)
			fmt.Fprintf(en.buf, "%s := packets.PutVarint(%s, int64(len(%s)))\n", n, x, b)
			fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
			fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, b))
		*/
		fmt.Fprintf(en.buf, "if err = packets.WriteString(ww, %s); err != nil { return err }", name)
	case "packets.VarInt", "packets.VarLong":
		x := en.T() // []byte for varint
		n := en.T() // num of varint bytes

		// TODO: Is making a new []byte every time efficient?
		fmt.Fprintf(en.buf, "%s := make([]byte, binary.MaxVarintLen64)\n", x)
		fmt.Fprintf(en.buf, "%s := packets.PutVarint(%s, int64(%s))\n", n, x, name)
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
	case "packets.Position", "packets.Angle":
		// TODO: Do we need to convert these to an appropriate type for binary.Write()?
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s)", Endianness, name))
	case "packets.UUID":
		fmt.Fprintf(en.buf, errWrap("binary.Write(ww, %s, %s[:])", Endianness, name))

	case "packets.Chat":
		// Covered by the 'as' switch above.
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
		fmt.Fprintf(en.buf, "// Unable to encode: %s (%s)", name, t)
	}

	fmt.Fprintf(en.buf, "\n")
}
