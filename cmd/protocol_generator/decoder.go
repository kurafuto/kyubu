package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
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
		tag := reflect.StructTag("")
		if field.Tag != nil {
			tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		}

		var c string
		if tag.Get("if") != "" && tag.Get("noreplace") != "true" {
			c = strings.Replace(tag.Get("if"), ".", name+".", -1)
		}

		if c != "" {
			fmt.Fprintf(de.buf, "if %s {\n", c)
		}

		for _, n := range field.Names {
			de.writeType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
		}

		if c != "" {
			fmt.Fprint(de.buf, "}\n")
		}
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
		lT := tag.Get("length") // length type
		if lT == "rest" {
			imports["io/ioutil"] = struct{}{}
			fmt.Fprintf(de.buf, "if %s, err = ioutil.ReadAll(rr); err != nil { return }\n", name)
			return
		}

		// length variable to read into
		lv := de.T()
		fmt.Fprintf(de.buf, "var %s %s\n", lv, lT)
		de.writeField(lT, lv, "")

		// TODO: Check for length > math.MaxInt16?

		imports["fmt"] = struct{}{}
		fmt.Fprintf(de.buf, `if %s < 0 {
				return fmt.Errorf("negative array size: %%d < 0", %s)
			}`, lv, lv)

		fmt.Fprintf(de.buf, "\n%s = make([]%s, %s)\n", name, e.Elt, lv)

		if i, ok := e.Elt.(*ast.Ident); ok && (i.Name == "byte" || i.Name == "uint8") {
			fmt.Fprintf(de.buf, "if _, err = rr.Read(%s); err != nil { return err }\n", name)
			return
		}

		iv := de.T()
		fmt.Fprintf(de.buf, "for %s := range %s {\n", iv, name)

		// Sneaky trick so we can use "unknown-ish" structs ([]T) and still get at their fields.
		sName := e.Elt.(*ast.Ident).Name
		if xx, ok := nonPackets[sName]; ok {
			de.writeType(xx, fmt.Sprintf("%s[%s]", name, iv), tag)
		} else if sName == "string" {
			de.writeField(sName, fmt.Sprintf("%s[%s]", name, iv), tag)
		} else {
			fmt.Fprintf(de.buf, "// Can't find supporting struct %s for %s.%s\n", sName, de.p.name, name)
		}

		fmt.Fprint(de.buf, "}\n")
	default:
		fmt.Fprintf(de.buf, "// Unable to decode: %s (type: %T)\n", name, e)
	}
}

func (de *Decoder) writeField(t, name string, tag reflect.StructTag) {
	as := tag.Get("as")
	if as != "" {
		switch as {
		case "json":
			imports["encoding/json"] = struct{}{}
			t := de.T()
			fmt.Fprintf(de.buf, `var %[1]s string
				if %[1]s, err = packets.ReadString(rr); err != nil { return err }
				if err = json.Unmarshal([]byte(%[1]s), &%[2]s); err != nil { return err }
			`, t, name)
		default:
			fmt.Fprintf(de.buf, "// Can't 'as' %s\n", as)
		}
		return
	}

	// TODO: For ints, unwrap binary.Read() trickery to reuse []byte tmp.
	switch t {
	case "bool":
		tmp := de.T() // byte value for bool
		fmt.Fprintf(de.buf, "var %s [1]byte\n", tmp)
		fmt.Fprintf(de.buf, "if _, err = rr.Read(%s[:1]); err != nil { return err }\n", tmp)
		fmt.Fprintf(de.buf, "%s = %s[0] == 0x01\n", name, tmp)
	case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
		fmt.Fprintf(de.buf, errWrap("binary.Read(rr, %s, %s)", Endianness, name))
	case "string":
		imports["errors"] = struct{}{}

		n := de.T() // num of varint bytes
		b := de.T() // []byte to read into
		x := de.T() // num of bytes read

		fmt.Fprintf(de.buf, "%s, err := packets.ReadVarint(rr)\n", n)
		fmt.Fprintf(de.buf, "if err != nil { return err}\n")
		fmt.Fprintf(de.buf, "%s := make([]byte, %s)\n", b, n)
		fmt.Fprintf(de.buf, "%s, err := rr.Read(%s)\n", x, b)
		fmt.Fprintf(de.buf, "if err != nil { return err } ")
		fmt.Fprintf(de.buf, "else if int64(%s) != %s {\n", x, n)
		// TODO: delta := n - n1; Read(b, delta) ...
		fmt.Fprintf(de.buf, `return errors.New("didn't read enough bytes for string")`+"\n")
		fmt.Fprintf(de.buf, "}\n")
		fmt.Fprintf(de.buf, "%s = string(%s)\n", name, b)
	case "packets.VarInt":
		n := de.T()
		fmt.Fprintf(de.buf, "%s, err := packets.ReadVarint(rr)\n", n)
		fmt.Fprintf(de.buf, "if err != nil { return err }\n")
		fmt.Fprintf(de.buf, "%s = packets.VarInt(%s)\n", name, n)
	case "packets.VarLong":
		n := de.T()
		fmt.Fprintf(de.buf, "%s, err := packets.ReadVarint(rr)\n", n)
		fmt.Fprintf(de.buf, "if err != nil { return err }\n")
		fmt.Fprintf(de.buf, "%s = packets.VarLong(%s)\n", name, n)
	case "packets.Position", "packets.Angle":
		// TODO: Do we need to take &%s so it reads into var?
		fmt.Fprintf(de.buf, errWrap("binary.Read(rr, %s, %s)", Endianness, name))
	case "packets.UUID":
		fmt.Fprintf(de.buf, errWrap("binary.Read(rr, %s, %s)", Endianness, name))

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
		fmt.Fprintf(de.buf, "// Unable to decode: %s (%s)", name, t)
	}

	fmt.Fprintf(de.buf, "\n")
}
